package imager

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

func Master(repositoryURLs string) Response {
	wg := &sync.WaitGroup{}
	ch := make(chan result)

	for _, repositoryURL := range strings.Split(repositoryURLs, ",") {
		repositories := repositories(repositoryURL)
		for _, repository := range repositories {
			wg.Add(1)
			go worker(repository, wg, ch)
		}
	}

	return reduce(wg, ch)
}

func worker(repository repository, wg *sync.WaitGroup, ch chan result) {
	if err := repository.Clone(); err != nil {
		log.WithField("instance", "worker").WithField("repository", repository).Error(err)
		wg.Done()
		return
	}

	err := filepath.Walk(repository.Path(), func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Only parse files with CMD or ENTRYPOINT as certain languages share Dockerfile semantics
		if strings.Contains(string(content), "CMD [") || strings.Contains(string(content), "ENTRYPOINT [") {
			data, err := parser.Parse(bytes.NewReader(content))
			if err != nil {
				return err
			}
			ch <- result{
				images:     images(data.AST),
				path:       strings.TrimPrefix(path, fmt.Sprintf("%s/", RepoDirectory)),
				repository: repository,
			}
		}

		return nil
	})

	if err != nil {
		ch <- result{error: err, repository: repository}
	}

	wg.Done()
}

func reduce(wg *sync.WaitGroup, ch chan result) map[string]map[string][]string {
	go func() {
		wg.Wait()
		close(ch)
	}()

	result := map[string]map[string][]string{}
	for data := range ch {
		if _, ok := result[data.repository.Key()]; !ok {
			result[data.repository.Key()] = map[string][]string{}
		}
		if _, ok := result[data.repository.Key()]["error"]; !ok {
			result[data.repository.Key()] = map[string][]string{"error": {}}
		}

		if data.error != nil {
			result[data.repository.Key()]["error"] = append(result[data.repository.Key()]["error"], data.error.Error())
		}

		result[data.repository.Key()][data.path] = data.images
	}

	return result
}
