package main

import (
	"encoding/json"
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

	r.HandleFunc("/", deck.BuildDeck(deck.TableTopSimulator, bulk, log))
	r.HandleFunc("/deck", deck.BuildDeck(deck.ScryfallEntry, bulk, log))

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

		deckList, err := mtgfail.ReadCardList(r, log)
		if err != nil {
			_, err = io.Copy(w, r)
			if err != nil {
				log.Error(
					"failed to parse data from request. Copying raw bytes through to response",
					"err", err,
				)
				http.Error(w, "", http.StatusOK)
				return
			}
		}
		defer r.Close()

		b, err := json.Marshal(deckList)
		if err != nil {

			log.Error(
				"Can't marshal deckfile",
				"err", err,
			)
			return

		}

		w.Header().Add("Content-Type", "application/json")

		_, err = fmt.Fprintf(w, "%s", b)
		if err != nil {
			log.Error(
				"Can't write deckfile",
				"err", err,
			)
			return

		}
		w.WriteHeader(200)

	},
	)

	if err = http.ListenAndServe(":8080", handlers.CORS(origins)(r)); err != nil {
		log.Error(
			"Server failure",
			"err", err,
		)
	}

}
