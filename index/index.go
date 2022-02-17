// basic-middleware.go
package main

import (
	"log"
	"net/http"
)

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		addCorsHeader(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println(r.URL.Path)
		f(w, r)
	}
}

func main() {
	http.HandleFunc("/containers/list", logging(containerList))
	http.HandleFunc("/containers/stats", logging(containerStats))
	http.HandleFunc("/containers/stop", logging(containerStop))
	http.HandleFunc("/containers/start", logging(containerStart))
	http.HandleFunc("/containers/create", logging(containerCreate))
	http.HandleFunc("/containers/restart", logging(containerRestart))
	http.HandleFunc("/containers/remove", logging(containerRemove))
	http.HandleFunc("/images/list", logging(imagesList))
	http.ListenAndServe(":8015", nil)
}
