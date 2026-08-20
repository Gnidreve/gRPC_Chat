package main

import (
	"encoding/json"
	"net/http"

	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

// onlineCountHandler serves the current subscriber count as JSON. It is the
// one deliberately public, unauthenticated endpoint on this server — the
// landing page (a different origin) polls it to show a live user count, and
// it exposes nothing beyond that single integer.
func onlineCountHandler(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(struct {
			Online int `json:"online"`
		}{Online: st.OnlineCount()})
	}
}
