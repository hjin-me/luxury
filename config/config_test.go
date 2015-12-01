package config

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	var (
		err error
	)
	type TestAppYaml struct {
		Env struct {
			Port int
			Tpl  string `ha:"path"`
		}
	}
	cfg := TestAppYaml{}
	// pwd log
	err = Load("test/path.yaml", &cfg)
	if err != nil {
		t.Error(err)
	}
	if cfg.Env.Port != 8088 {
		t.Error("cfg error", cfg)
	}
	x, _ := os.Getwd()
	if cfg.Env.Tpl != filepath.Join(x, "test/output") {
		t.Error(cfg)
	}
	t.Skip()

	type TestYaml struct {
		Abs string `ha:"path"`
	}
	tcfg := TestYaml{}
	// absolute log
	filename, _ := filepath.Abs("test/abs.yaml")
	err = Load(filename, &tcfg)
	if err != nil {
		t.Error(err)
	}

	// not exists log
	err = Load("test", &cfg)
	if err == nil {
		t.Error("should cause an error")
	}

}
func TestReflect(t *testing.T) {
	t.Skip()

	type TestAppYaml struct {
		Env struct {
			Port string `ha:"path"`
		}
		T string `ha:"absolute"`
	}

	data := TestAppYaml{}
	data.Env.Port = "test"
	data.T = "1234"

	scan(reflect.ValueOf(&data))
	t.Log(data)

	data = TestAppYaml{}
	data.Env.Port = "zxc"
	data.T = "123"
	scan(reflect.ValueOf(data))

	t.Log(data)

	t.Error("end")

}
func scan(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	log.Println(v.Kind())
	switch v.Kind() {
	case reflect.Struct:
		n := v.NumField()
		t := v.Type()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.Struct:
				scan(f)
			case reflect.String:
				tf := t.Field(i)
				if s := tf.Tag.Get("ha"); s != "" && f.CanSet() {
					f.SetString("xxxxxx:" + f.String())
				}
			}
		}
	}

}
