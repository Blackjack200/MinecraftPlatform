package main

import (
	"Blackjack200/MinecraftPlatform/api"
	"Blackjack200/MinecraftPlatform/storage"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//goland:noinspection GoUnhandledErrorResult
func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	logrus.Info("Start Platform")
	if err := storage.Initialize(); err != nil {
		logrus.Fatal(err)
	}

	c, _ := json.Marshal(storage.Config.Servers)
	api.InitializeList(string(c))
	api.InitializeQuery()
	templateBox := packr.NewBox("./templates/")
	index, _ := templateBox.FindString("index.html")
	entry, _ := templateBox.FindString("entry.html")
	e := base64.StdEncoding.EncodeToString([]byte(entry))
	t, _ := template.New("index").Parse(index)

	staticBox := packr.NewBox("./static/")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		parts := strings.Split(strings.ToLower(request.RequestURI), "/")
		part := parts[len(parts)-1]
		if strings.HasPrefix(part, "query") {
			api.HandleQuery(writer, request)
			return
		}
		if strings.HasPrefix(part, "list") {
			api.HandleServerList(writer, request)
			return
		}

		if len(parts) >= 2 && parts[1] == "static" {
			ctx, err := staticBox.FindString(strings.TrimPrefix(part, "static"))
			if err != nil {
				fmt.Fprint(writer, err)
				return
			}
			fmt.Fprint(writer, ctx)
			return
		}
		if err := t.Execute(writer, struct {
			Title string
			Entry string
		}{
			Title: "ServerPlatform",
			Entry: e,
		}); err != nil {
			logrus.Fatal(err)
		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sig
		logrus.Info("Save Config")
		if err, err2 := storage.Save(); err != nil {
			logrus.Fatal(err)
		} else if err2 != nil {
			logrus.Fatal(err2)
		}
		os.Exit(0)
	}()
	errs := make(chan error)
	go func() {
		if err := http.ListenAndServe(storage.Config.Bind, nil); err != nil {
			errs <- err
			logrus.Fatal(err)
		} else {
		}
	}()
	logrus.Info("Bind: " + storage.Config.Bind)
	<-errs
}
