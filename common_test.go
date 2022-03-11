package imager

import (
	"bufio"
	"bytes"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"
)

func TestImages(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/Dockerfile")
	if err != nil {
		t.Error(err)
	}

	data, err := parser.Parse(bytes.NewReader(content))
	if err != nil {
		t.Error(err)
	}

	expected := []string{"golang:1.13", "alpine"}
	got := images(data.AST)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %#+v, got %#+v", expected, got)
	}
}

func TestSplitEntryIntoRepository(t *testing.T) {
	entry := "https://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404"
	got, err := splitEntryIntoRepository(entry)
	if err != nil {
		t.Error(err)
	}

	expectedURL, err := url.Parse("https://github.com/app-sre/container-images.git")
	if err != nil {
		t.Error(err)
	}
	expected := &repository{url: *expectedURL, commitSHA: "c260deaf135fc0efaab365ea234a5b86b3ead404"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %#+v, got %#+v", expected, got)
	}
}

func TestRepositories(t *testing.T) {
	inputURL := "https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt"
	got := repositories(inputURL)

	var expected []repository
	content, err := ioutil.ReadFile("fixtures/input.txt")
	if err != nil {
		t.Error(err)
	}

	buffer := bufio.NewScanner(bytes.NewReader(content))
	for buffer.Scan() {
		repo, err := splitEntryIntoRepository(buffer.Text())
		if err != nil {
			t.Error(err)
		}
		expected = append(expected, *repo)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %#+v, got %#+v", expected, got)
	}
}
