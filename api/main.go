package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func main() {
	// Ayyy pprof!
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "大丈夫")
		})
		logErr(http.ListenAndServe(":8000", nil))
	}()

	// HTTP to HTTPS redirect, with sneaky cert renewal support
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/.well-known/", http.StripPrefix("/.well-known", http.FileServer(http.Dir(".well-known"))))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Scheme = "https"
			u.Host = r.Host
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		})
		logErr(http.ListenAndServe(":800", mux))
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "大丈夫")
	})
	logErr(http.ListenAndServe(":80", HSTS(CORS(mux))))
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", strings.ToUpper(r.Header.Get("Access-Control-Request-Method")))
		headers.Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	})
}

func HSTS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Strict-Transport-Security", "max-age=315360000; preload")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-XSS-Protection", "1; mode=block")
		headers.Set("Content-Security-Policy", "default-src 'self'")

		h.ServeHTTP(w, r)
	})
}
