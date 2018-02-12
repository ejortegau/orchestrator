package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/openark/golib/log"

	"github.com/github/orchestrator/go/config"
	"github.com/github/orchestrator/go/raft"
	"github.com/go-martini/martini"
)

var reverseProxy = func(w http.ResponseWriter, r *http.Request, c martini.Context) {
	if !orcraft.IsRaftEnabled() {
		return
	}
	if orcraft.IsLeader() {
		return
	}
	leader := orcraft.GetLeader()
	if leader == "" {
		log.Errorf("raft reverse-proxy: leader is unknown")
		return
	}
	hostPort := strings.Split(leader, ":")
	leaderHost := hostPort[0]

	hostPort = strings.Split(config.Config.ListenAddress, ":")
	port := hostPort[1]
	leaderURI := fmt.Sprintf("http://%s:%s", leaderHost, port)
	url, err := url.Parse(leaderURI)
	if err != nil {
		log.Errore(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	log.Debugf("................reverse proxy to %s", leaderURI)
	proxy.ServeHTTP(w, r)
}

// All returns a Handler that adds gzip compression to all requests
func RaftReverseProxy() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, c martini.Context) {
		reverseProxy(w, r, c)
	}
}
