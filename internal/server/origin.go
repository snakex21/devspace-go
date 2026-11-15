package server

import (
	"net"
	"net/http"
	"strings"

	"github.com/snakex21/devspace-go/internal/workspace"
)

// forwardingHeaders are added by Cloudflare, Pinggy and other reverse
// proxies. A client on the internet cannot remove the ones the proxy adds.
var forwardingHeaders = []string{
	"Cf-Connecting-Ip", "Cf-Ray", "Cf-Ipcountry", "Cdn-Loop",
	"X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto",
	"X-Real-Ip", "True-Client-Ip", "Forwarded", "Via",
}

// requestOrigin classifies a request as local only when every signal says
// so: the TCP peer is a loopback address, no proxy forwarding header is
// present, and the Host header names this machine. A tunnel connects from
// loopback too, but it forwards the public hostname and adds forwarding
// headers. Anything uncertain is remote, the stricter side; a local client
// that sends such headers only restricts itself.
func requestOrigin(r *http.Request, listenHost string) workspace.Origin {
	if r == nil {
		return workspace.OriginRemote
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return workspace.OriginRemote
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		return workspace.OriginRemote
	}
	for _, name := range forwardingHeaders {
		if _, ok := r.Header[http.CanonicalHeaderKey(name)]; ok {
			return workspace.OriginRemote
		}
	}
	if !isLocalHostHeader(r.Host, listenHost) {
		return workspace.OriginRemote
	}
	return workspace.OriginLocal
}

func isLocalHostHeader(hostHeader, listenHost string) bool {
	name := hostHeader
	if h, _, err := net.SplitHostPort(hostHeader); err == nil {
		name = h
	}
	name = strings.ToLower(strings.Trim(name, "[]"))
	if name == "localhost" || strings.HasSuffix(name, ".localhost") {
		return true
	}
	if ip := net.ParseIP(name); ip != nil {
		return ip.IsLoopback()
	}
	listenHost = strings.ToLower(strings.Trim(listenHost, "[]"))
	return name != "" && name == listenHost && listenHost != "0.0.0.0" && listenHost != "::"
}
