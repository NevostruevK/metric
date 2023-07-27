package handlers

import (
	"fmt"
	"net"
	"net/http"
)

func IPCheckHandler(next http.Handler, ipNet *net.IPNet) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sIP := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(sIP)
		if ip == nil {
			msg := fmt.Sprintf(`can't parse IP from header "X-Real-IP" %s`, sIP)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		if !ipNet.Contains(ip) {
			msg := fmt.Sprintf(`IP %s is forbidden for using this application`, ip)
			http.Error(w, msg, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
