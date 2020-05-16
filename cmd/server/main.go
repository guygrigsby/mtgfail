package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/guygrigsby/mtgfail"
	"github.com/guygrigsby/mtgfail/deck"
	"github.com/inconshreveable/log15"
)

var (
	bulk string
)

func main() {

	flag.StringVar(&bulk, "bulk", "./scryfall-default-cards.json", "The bulk data download")
	flag.Parse()

	log := log15.New()

	bulk, err := mtgfail.ReadBulk(bulk, log)
	if err != nil {
		log.Error(
			"Can't parse bulk data",
			"err", err,
		)
		return
	}
	log.Info(
		"Read bulk",
		"entries", len(bulk),
	)

	r := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r.HandleFunc("/", deck.BuildDeck(bulk, log))
	r.HandleFunc("/site", func(w http.ResponseWriter, req *http.Request) {
		log.Debug(
			"Request",
			"req", fmt.Sprintf("%+v", req),
		)

		r, err, status := deck.FetchDeck(req, log)
		if status > 299 {
			log.Error(
				"failed to fetch deck",
				"err", err,
			)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(w, r)
		if err != nil {
			log.Error(
				"failed to copy data from request",
				"err", err,
			)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		r.Close()

	},
	)

	if err = http.ListenAndServe(":8080", handlers.CORS(origins)(r)); err != nil {
		log.Error(
			"Server failure",
			"err", err,
		)
	}

}
