package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	var (
		err error
	)
	t.Log(os.Getwd())
	// pwd log
	f, err := Load("test/path.yaml")
	if err != nil {
		t.Error(err)
	}
	if len(f.Data) == 0 {
		t.Error(f, "length is 0")
	}
	// absolute log
	filename, _ := filepath.Abs("test/abs.yaml")
	t.Log(filename)
	f, err = Load(filename)
	if err != nil {
		t.Error(err)
	}
	if len(f.Data) > 0 {
		t.Error(f, "length should be 0")
	}

	// not exists log
	f, err = Load("test")
	if err == nil {
		t.Error("should cause an error")
	}
	if len(f.Data) > 0 {
		t.Error(f, "length should be 0")
	}
}
