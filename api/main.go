package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/fugiman/tyrantbot/twitch"
)

type Config struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var config Config
var db = dynamodb.New(session.New(aws.NewConfig().WithRegion("us-west-2")))

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	twitch.API.ClientId = Config.ClientId
	twitch.API.ClientSecret = Config.ClientSecret
}

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

	// Define our routes
	mux := web.New()
	mux.Use(middleware.EnvInit)
	mux.Use(HSTS)
	mux.Use(CORS)
	mux.Use(Sessions)
	mux.Use(context.ClearHandler)

	mux.Get("/", home)
	mux.Get("/callback", login)

	// Run the server gracefully
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
	graceful.DoubleKickWindow(2 * time.Second)
	listener := bind.Socket(":80")
	mux.Compile()
	log.Println("Starting Goji on", listener.Addr())
	graceful.HandleSignals()
	bind.Ready()
	graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
	graceful.PostHook(func() { log.Printf("Goji stopped") })
	logErr(graceful.Serve(listener, mux))
	graceful.Wait()
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HSTS(c *web.C, h http.Handler) http.Handler {
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

func CORS(c *web.C, h http.Handler) http.Handler {
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

func Sessions(c *web.C, h http.Handler) http.Handler {
	store := sessions.NewCookieStore([]byte("dickbutt"))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, e := store.Get(r, "session")
		if e != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error: %s", e)
			return
		}
		if ip, ok := s.Values["ip"]; ok {
			if ip.(string) != r.RemoteAddr {
				delete(s.Values, "ip")
				w.WriteHeader(401)
				fmt.Fprintf(w, "Session hijacking detected")
				return
			}
		} else {
			s.Values["ip"] = r.RemoteAddr
		}

		c.Env["session"] = s
		h.ServeHTTP(w, r)
	})
}
