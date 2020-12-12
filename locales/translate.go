package locales

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var locales = make(map[string]*map[string]string)

func loadDictionary(path string) (*map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return nil, err2
	}

	var en map[string]string
	json.Unmarshal(byteValue, &en)

	return &en, nil
}

func PreloadLocales(dictAssets string) error {
	err := filepath.Walk(dictAssets, func(path string, info os.FileInfo, err error) error {
		hoh := filepath.Ext(path)
		if info.IsDir() || hoh != ".json" {
			return nil
		}

		dict, err2 := loadDictionary(path)
		if err2 != nil {
			return err2
		}

		code := strings.Split(info.Name(), ".")[0]
		locales[code] = dict
		return nil
	})
	return err
}

func Translate(lng string, key string) string {
	dict := locales[lng]
	if dict == nil {
		dict = locales["en"]
	}

	val := (*dict)[key]
	if val == "" {
		val = (*locales["en"])[key]
	}

	return val
}
