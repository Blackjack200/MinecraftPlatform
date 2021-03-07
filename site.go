package main

import (
	"Blackjack200/MinecraftPlatform/config"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/sandertv/gophertunnel/query"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var listLock sync.Mutex
var Cache string

func main() {
	if err := config.Initialize(); err != nil {
		panic(err)
	}
	c, _ := json.Marshal(config.Cfg.Servers)
	Cache = string(c)
	listLock = sync.Mutex{}
	box := packr.NewBox("./templates")
	cb, _ := box.FindString("index.html")
	t, _ := template.New("index").Parse(cb)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		parts := strings.Split(strings.ToLower(request.RequestURI), "/")
		part := parts[len(parts)-1]
		if strings.HasPrefix(part, "query") {
			HandleQuery(writer, request)
			return
		}

		if strings.HasPrefix(part, "list") {
			HandleServerList(writer, request)
			return
		}
		if err := t.Execute(writer, struct {
			Title string
		}{Title: "This is title"}); err != nil {
			panic(err)
		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sig
		println("saved")
		if err := config.Save(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	err := http.ListenAndServe(config.Cfg.Bind, nil)
	if err != nil {
		panic(err)
	}
}

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()
	err1 := r.ParseForm()
	if err1 != nil {
		_, _ = fmt.Fprint(w, "{\"error\":\"invalid request\"}")
		return
	}
	host, port := r.Form.Get("host"), r.Form.Get("port")
	var result map[string]string
	d, er := query.Do(host + ":" + port)
	if er != nil {
		result = nil
	} else {
		result = d
	}

	if result == nil {
		_, _ = fmt.Fprint(w, "{\"error\":\"timeout\"}")
		return
	}

	req, _ := json.Marshal(result)
	_, _ = fmt.Fprint(w, string(req))
}

func HandleServerList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		listLock.Unlock()
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	listLock.Lock()
	err := r.ParseForm()

	if err != nil {
		_, _ = fmt.Fprint(w, "{\"error\":\"invalid request\"}")
		return
	}

	host, port := r.Form.Get("host"), r.Form.Get("port")
	if r.Form.Get("get") != "" {
		_, _ = fmt.Fprint(w, Cache)
		return
	}
	record := strings.TrimSpace(host) + ":" + strings.TrimSpace(port)
	_, ex := config.Cfg.Servers[record]
	if ex {
		_, _ = fmt.Fprint(w, "{\"error\":\"existed\"}")
		return
	} else {
		re, er := GetHostByName(host)
		if er == nil {
			config.Cfg.Servers[record] = re
			c, _ := json.Marshal(config.Cfg.Servers)
			Cache = string(c)
			_, _ = fmt.Fprint(w, "{\"success\":\"_\"}")
		} else {
			_, _ = fmt.Fprint(w, "{\"error\":\"host\"}")
		}
	}
}

func GetHostByName(host string) (string, error) {
	addr, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, ip := range addr {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}
	return "", &net.DNSError{}
}
