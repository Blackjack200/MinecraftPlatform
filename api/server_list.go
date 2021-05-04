package api

import (
	"Blackjack200/MinecraftPlatform/storage"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var mu sync.Mutex
var listCache string

func InitializeList(list string) {
	listCache = list
	mu = sync.Mutex{}
}

func HandleServerList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	if ParseRequest(w, r) {
		return
	}

	host, port := r.Form.Get("host"), r.Form.Get("port")
	if r.Form.Get("get") != "" {
		fmt.Fprint(w, listCache)
		return
	}

	record := FormatAddr(host, port)
	mu.Lock()
	defer mu.Unlock()
	if _, contains := storage.Config.Servers[record]; contains {
		fmt.Fprint(w, "{\"error\":\"existed\"}")
		return
	} else {
		if ret, err := GetHostByName(host); err == nil {
			storage.Config.Servers[record] = ret
			logrus.Info("ADD: " + record)
			refreshCacheList()
			fmt.Fprint(w, "{\"success\":\"_\"}")
		} else {
			fmt.Fprint(w, "{\"error\":\"host\"}")
		}
	}
}

func refreshCacheList() {
	c, _ := json.Marshal(storage.Config.Servers)
	listCache = string(c)
}
