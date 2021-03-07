package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Bind    string            `json:"local"`
	Servers map[string]string `json:"server-list"`
}

var (
	Cfg        *config
	File       string
	WorkingDir string
)

func Initialize() error {
	WorkingDir, _ = os.Getwd()
	WorkingDir, _ = filepath.Abs(WorkingDir)
	File = WorkingDir + "/config.json"

	Cfg = &config{}
	Cfg.Bind = ":888"
	Cfg.Servers = make(map[string]string)

	if _, s := os.Stat(File); os.IsNotExist(s) {
		bytes, err := json.MarshalIndent(Cfg, "", "	")
		if err != nil {
			return err
		}
		_ = ioutil.WriteFile(File, bytes, 0777)
	} else {
		con, _ := ioutil.ReadFile(File)
		err := json.Unmarshal(con, Cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func Save() error {
	bytes, err := json.MarshalIndent(Cfg, "", "	")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(File, bytes, 0777)
	return err
}
