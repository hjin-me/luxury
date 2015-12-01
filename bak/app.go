package main

import (
	"config"
	"db"
	"flag"
	"install"
	"logex"
	"net/http"
	"path/filepath"
	"route"
	"theme"
	"time"
	"workflow"

	"golang.org/x/net/context"
)

type DbConf struct {
	DSN string `yaml:"dsn"`
}

var (
	isInstall = flag.Bool("install", false, "install Haruhi")
)

func main() {
	logex.SetLogLevel(logex.DEBUG)
	flag.Parse()
	// init config loader
	if len(flag.Args()) > 0 {
		baseDir := flag.Args()[0]
		if !filepath.IsAbs(baseDir) {
			var err error
			baseDir, err = filepath.Abs(baseDir)
			if err != nil {
				logex.Fatal(err)
			}
		}
		config.SetBaseDir(baseDir)
		logex.Trace("base dir is ", baseDir)
	}

	// connect database
	dsn := DbConf{}
	err := config.Load("db.yaml", &dsn)
	if err != nil {
		logex.Fatal(err, "sql dsn not ok")
	}
	err = db.Create(dsn.DSN)
	if err != nil {
		logex.Fatal(err)
	}
	defer db.Close()

	switch *isInstall {
	case true:
		install.Install()
		return
	default:

		// start server
		ctx := getMux()
		go func() {
			err := http.ListenAndServe(":8088", ctx) //设置监听的端口
			if err != nil {
				ctx.Cancel()
				logex.Trace(err)
			}

		}()
		<-ctx.Done()

	}

}

func getMux() *MuxContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &MuxContext{ctx, cancel}
}

type MuxContext struct {
	context.Context
	Cancel func()
}

func (p *MuxContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	statusCode := http.StatusOK
	var err interface{}
	defer func() {
		endTime := time.Now()
		if err == nil {
			err = ""
		}
		logex.Noticef("%d %0.3f %s %s [%s] [%s]", statusCode, float32(endTime.Sub(startTime))/float32(time.Second), r.Method, r.URL.Path, r.URL.RawQuery, err)
	}()

	wf := workflow.New()
	wf.Context = context.WithValue(context.WithValue(wf.Context, "req", r), "res", w)

	workflow.NewAgent(wf.Context, func() {
	}, "start")
	workflow.NewAgent(wf.Context, logex.Agent, "logger")
	workflow.NewAgent(wf.Context, route.File, "static")
	workflow.NewAgent(wf.Context, route.Domains, "domain")
	workflow.NewAgent(wf.Context, route.Url, "url")
	workflow.NewAgent(wf.Context, theme.Agent, "theme")

	err = wf.Load("workflow.yaml")
	if err != nil {
		return
	}
	err = wf.Handle("start")
	if err != nil {
		logex.Warning(err)
	}
}
