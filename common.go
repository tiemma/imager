package imager

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

const (
	RepoDirectory = "repositories"
	from          = "from"
)

func images(node *parser.Node) []string {
	var images []string
	for _, child := range node.Children {
		if strings.ToLower(child.Value) == from {
			images = append(images, child.Next.Value)
		}
	}

	return images
}

func splitEntryIntoRepository(entry string) (*repository, error) {
	values := strings.Split(entry, " ")
	if len(values) != 2 {
		return nil, fmt.Errorf("entry must be of length 2: %s", entry)

	}

	gitURL, err := url.Parse(values[0])
	if err != nil {
		return nil, err
	}

	return &repository{url: *gitURL, commitSHA: values[1]}, nil
}

func openRemoteFile(inputURL string) io.ReadCloser {
	resp, err := http.Get(inputURL)
	if err != nil {
		log.WithField("instance", "openRemoteFile").WithField("url", inputURL).Fatal(err)
	}

	return resp.Body
}

func repositories(inputURLs string) []repository {
	reader := openRemoteFile(inputURLs)
	defer func() {
		if err := reader.Close(); err != nil {
			log.WithField("instance", "repositories").WithField("url", inputURLs).Fatal(err)
		}
	}()

	var repositories []repository
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		repositoryEntry := scanner.Text()
		repository, err := splitEntryIntoRepository(repositoryEntry)
		if err != nil {
			log.WithField("instance", "repositories").WithField("url", inputURLs).Errorf("Skipping %s: %s", repositoryEntry, err)
			continue
		}
		repositories = append(repositories, *repository)
	}

	return repositories
}
