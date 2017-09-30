package engine

import (
	"io/ioutil"
	"os"
	"testing"
)

var mockConfig = `
{
	"version": 1
}
`

func TestParse(t *testing.T) {
	_, err := ParseString(mockConfig)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = ParseString("[]")
	if err == nil {
		t.Errorf("Want parse error, got nil")
	}
}

func TestParseFile(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "drone")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(mockConfig)
	f.Close()

	_, err = ParseFile(f.Name())
	if err != nil {
		t.Error(err)
		return
	}

	_, err = ParseFile("/tmp/this/path/does/not/exist")
	if err == nil {
		t.Errorf("Want parse error, got nil")
	}
}
