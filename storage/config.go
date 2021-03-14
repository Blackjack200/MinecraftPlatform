package storage

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Bind    string            `json:"local"`
	DBPath  string            `json:"db-path"`
	Servers map[string]string `json:"server-list"`
}

var (
	Config     *config
	File       string
	WorkingDir string
	DB         *leveldb.DB
)

func Initialize() error {
	WorkingDir, _ = os.Getwd()
	WorkingDir, _ = filepath.Abs(WorkingDir)
	File = WorkingDir + "/config.json"

	Config = &config{}
	Config.Bind = ":888"
	Config.DBPath = "db"
	Config.Servers = make(map[string]string)

	if _, s := os.Stat(File); os.IsNotExist(s) {
		bytes, err := json.MarshalIndent(Config, "", "	")
		if err != nil {
			return err
		}
		_ = ioutil.WriteFile(File, bytes, 0777)
	} else {
		con, _ := ioutil.ReadFile(File)
		err := json.Unmarshal(con, Config)
		if err != nil {
			return err
		}
	}
	DB, _ = leveldb.OpenFile(WorkingDir+"/"+Config.DBPath, nil)
	if DB == nil {
		return errors.New("LevelDB Open Failed")
	}
	logrus.Info("DB: " + WorkingDir + "/" + Config.DBPath)
	return nil
}

func Save() (error, error) {
	err2 := DB.Close()
	bytes, err := json.MarshalIndent(Config, "", "	")
	if err != nil {
		return err, err2
	}
	err = ioutil.WriteFile(File, bytes, 0777)
	return err, err2
}
