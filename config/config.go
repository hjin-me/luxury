package config

import (
	"file"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

func configScan(v reflect.Value, base string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		n := v.NumField()
		t := v.Type()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.Struct:
				configScan(f, base)
			case reflect.String:
				tf := t.Field(i)
				if s := tf.Tag.Get("ha"); s == "path" && f.CanSet() {
					if !filepath.IsAbs(f.String()) {
						f.SetString(filepath.Join(base, f.String()))
					}
				}
			}
		}
	}
}
func configUnmarshal(bf []byte, data interface{}, filename string) (err error) {
	err = yaml.Unmarshal(bf, data)
	if err != nil {
		return
	}
	rv := reflect.ValueOf(data)
	configScan(rv, filepath.Dir(filename))
	return
}

func SetBaseDir(dir string) {
	file.SetBaseDir(dir)
}

func Load(filename string, data interface{}) (err error) {
	file, err := file.Load(filename)
	if err != nil {
		return
	}
	err = configUnmarshal(file.Data, data, file.Name)
	return
}
