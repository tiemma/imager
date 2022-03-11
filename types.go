package imager

import (
	"fmt"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"net/url"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type repository struct {
	url       url.URL
	commitSHA string
}

type result struct {
	images []string
	path   string
	error  error
	repository
}

type Response map[string]map[string][]string

func (r repository) Key() string {
	return fmt.Sprintf("%s:%s", r.url.String(), r.commitSHA)
}

func (r repository) Path() string {
	urlPaths := strings.Split(r.url.String(), "/")
	repositoryName := strings.Split(urlPaths[len(urlPaths)-1], ".")[0]

	return fmt.Sprintf("%s/%s", RepoDirectory, repositoryName)
}

func (r repository) CloneOptions() *git.CloneOptions {
	return &git.CloneOptions{
		URL: r.url.String(),
	}
}

func (r repository) Clone() error {
	repo, err := git.Clone(
		memory.NewStorage(),
		filesystem.NewStorage(osfs.New(r.Path()), cache.NewObjectLRUDefault()).Filesystem(),
		r.CloneOptions(),
	)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := wt.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(r.commitSHA)}); err != nil {
		return err
	}

	return nil
}
