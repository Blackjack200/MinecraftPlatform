package api

type serverStatus struct {
	Uptime  []uptime `json:"uptime"`
	Maximum []uptime `json:"maximum"`
}

type uptime struct {
}
