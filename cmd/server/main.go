package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
)

func BuildDeck(cache mtgfail.Bulk, log log15.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// GCP load balancer health checks are garbage. Somehow, they always end up at '/'
			// This was I don' spend hours softing out why my pods are unhealthy. TODO fix it right
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var (
			err error
		)

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(
				"Cannot read body",
				"err", err,
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Debug(
			"request",
			"content", fmt.Sprintf("%+s", content),
		)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(content))

		deckList := make(map[string]int)

		switch r.Header.Get("Content-Type") {
		case "application/json":
			//TODO
		case "application/x-www-form-urlencoded":
			//TODO
		case "text/plain":
			fallthrough
		default:
			deckList, err = mtgfail.ReadCardList(r.Body, log)
			if err != nil {
				log.Error(
					"Can't read cardfile",
					"err", err,
				)
				http.Error(w, "Can't read card list", http.StatusBadRequest)
				return
			}

		}

		deck, err := mtgfail.BuildDeck(cache, deckList, log)
		if err != nil {
			log.Error(
				"Cannot build deck",
				"err", err,
			)
			return
		}

		b, err := json.Marshal(deck)
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
	}
}

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

	http.HandleFunc("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	http.HandleFunc("/", BuildDeck(bulk, log))
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Error(
			"Server failure",
			"err", err,
		)
	}

}
