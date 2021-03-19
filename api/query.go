package api

import (
	"encoding/json"
	"fmt"
	"github.com/sandertv/gophertunnel/query"
	"net/http"
	"sync"
	"time"
)

var queryCache sync.Map

func InitializeQuery() {
	queryCache = sync.Map{}
	go func() {
		for {
			queryCache = sync.Map{}
			time.Sleep(time.Second * 5)
		}
	}()
}

//goland:noinspection GoUnhandledErrorResult
func HandleQuery(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	if ParseRequest(w, r) {
		return
	}

	addr := FormatAddr(r.Form.Get("host"), r.Form.Get("port"))
	var result string

	if data, contains := queryCache.Load(addr); contains {
		cache, _ := data.(string)
		result = cache
	} else {
		if ret, err := query.Do(addr); err != nil {
			fmt.Fprint(w, "{\"error\":\"timeout\"}")
			return
		} else {
			ctx, _ := json.Marshal(ret)
			result = string(ctx)
			queryCache.Store(addr, string(ctx))
		}
	}

	_, _ = fmt.Fprint(w, result)
}
