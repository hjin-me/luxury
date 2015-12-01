package route

import (
	"config"
	"logex"
	"net/http"
	"strings"
	"sync"
)

type StaticConfig struct {
	Dir string `ha:"path"`
}

var (
	once     sync.Once
	fsHandle http.Handler
)

func File(log logex.Logger, r *http.Request, w http.ResponseWriter) (string, error) {
	var (
		err error
	)
	log.Debug("url path", r.URL.Path)
	if strings.Index(r.URL.Path, "/favicon.ico") == 0 {
		log.Debug("match favicon rule")
		return "end", nil
	}
	if strings.Index(r.URL.Path, "/static") != 0 {
		return "through", nil
	}

	once.Do(func() {
		cfg := StaticConfig{}
		err = config.Load("static.yaml", &cfg)
		log.Debug(cfg.Dir)
		if err != nil {
			return
		}
		fsHandle = http.FileServer(http.Dir(cfg.Dir))
	})
	if err != nil {
		return "", err
	}

	fsHandle.ServeHTTP(w, r)
	return "end", nil
}
