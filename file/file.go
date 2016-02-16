package file

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/hjin-me/luxury/logex"

	"golang.org/x/net/context"
)

var (
	baseDir string
	envDir  string = os.Getenv("HPATH")
)

func SetBaseDir(dir string) {
	baseDir = dir
}

var (
	fileCacheMap = make(map[string]*cCache)
)

type cCache struct {
	once sync.Once
	File File
}
type File struct {
	Data []byte
	Name string
}

func Load(filename string) (outerFile File, outerErr error) {
	cache, ok := fileCacheMap[filename]
	if !ok {
		cache = &cCache{}
		fileCacheMap[filename] = cache
	} else {
		outerFile = cache.File
	}
	cache.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		pwd, _ := os.Getwd()

		select {
		case filename = <-absFilepath(ctx, filename):
		case filename = <-relativeFilepath(ctx, pwd, filename):
		case filename = <-relativeFilepath(ctx, baseDir, filename):
		case filename = <-relativeFilepath(ctx, envDir, filename):
		case <-ctx.Done():
			outerErr = errors.New("cant find file [" + filename + "]")
			return
		}

		logex.Debug(filename)
		f, err := os.Open(filename)
		if err != nil {
			outerErr = err
			return
		}
		defer f.Close()

		outerFile.Name = filename
		outerFile.Data, err = ioutil.ReadAll(f)
		if err != nil {
			outerErr = err
			return
		}
		cache.File = outerFile
	})
	return
}
