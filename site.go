package main

import (
	"Blackjack200/MinecraftPlatform/config"
	"encoding/json"
	"fmt"
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
	config.Initialize()
	c, _ := json.Marshal(config.Cfg.Servers)
	Cache = string(c)
	listLock = sync.Mutex{}
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
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(writer, struct {
			Title string
		}{Title: "This is title"})
	})

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigs
		config.Save()
		println("saved")
	}()

	err := http.ListenAndServe(":666", nil)
	if err != nil {
		panic(err)
	}
}

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
		fmt.Fprint(w, "{\"error\":\"timeout\"}")
		return
	}

	req, _ := json.Marshal(result)
	_, _ = fmt.Fprint(w, string(req))
}

func HandleServerList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
		listLock.Unlock()
	}()
	listLock.Lock()
	err1 := r.ParseForm()
	if err1 != nil {
		_, _ = fmt.Fprint(w, "{\"error\":\"invalid request\"}")
		return
	}
	host, port := r.Form.Get("host"), r.Form.Get("port")
	if r.Form.Get("get") != "" {
		fmt.Fprint(w, Cache)
		return
	}
	record := strings.TrimSpace(host) + ":" + strings.TrimSpace(port)
	_, ex := config.Cfg.Servers[record]
	if ex {
		fmt.Fprint(w, "{\"error\":\"existed\"}")
		return
	} else {
		re, er := GetHostByName(host)
		if er == nil {
			config.Cfg.Servers[record] = re
			c, _ := json.Marshal(config.Cfg.Servers)
			Cache = string(c)
			fmt.Fprint(w, "{\"success\":\"_\"}")
		} else {
			fmt.Fprint(w, "{\"error\":\"host\"}")
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

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
