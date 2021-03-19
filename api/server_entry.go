package api

import (
	"time"
)

type serverStatus struct {
	Uptime  []uptime `json:"uptime"`
	Maximum []uptime `json:"maximum"`
}

type uptime struct {
	Time time.Time `json:"time"`
	Data []string  `json:"data"`
}
