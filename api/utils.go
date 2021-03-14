package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

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

func FormatAddr(host string, port string) string {
	return strings.TrimSpace(host) + ":" + strings.TrimSpace(port)
}

func ParseRequest(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprint(w, "{\"error\":\"invalid request\"}")
		return true
	}
	return false
}
