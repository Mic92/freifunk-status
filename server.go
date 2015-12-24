package main

import (
	"encoding/json"
	"net/http"

	"github.com/Mic92/freifunk-status/sysinfo"
)

func serveSysinfo(w http.ResponseWriter, r *http.Request) {
	info := sysinfo.New()

	js, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func main() {
	http.HandleFunc("/sysinfo-json.cgi", serveSysinfo)
	serveSingle("/", "static/index.html")
	http.ListenAndServe(":80", nil)
}
