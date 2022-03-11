package imager

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

const (
	inputURL = "https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt"
)

func TestMaster(t *testing.T) {
	output, err := ioutil.ReadFile("fixtures/expected.json")
	if err != nil {
		t.Error(err)
	}

	var data Response
	if err := json.Unmarshal(output, &data); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(Master(inputURL), data) {
		t.Error("Structs are not equal")
	}
}
