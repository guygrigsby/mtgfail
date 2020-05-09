package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/avast/retry-go"
	"github.com/guygrigsby/mtgfail"
	tts "github.com/guygrigsby/mtgfail/pkg/tabletopsimulator"
	"github.com/inconshreveable/log15"
)

// BuildDeck ...
func BuildDeck(cache mtgfail.Bulk, log log15.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if (*req).Method == "OPTIONS" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			//w.Header().Add("Access-Control-Allow-Origin", "https://mtg.fail, https://api.mtg.fail")
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Cache-Control")
			return
		}
		var (
			err error
			r   io.ReadCloser
		)
		if req.Method == http.MethodGet || req.Method == http.MethodPost {
			q := req.URL.Query()
			if len(q) == 0 {
				// GCP load balancer health checks are garbage. Somehow, they always end up at '/'
				// This was I don' spend hours softing out why my pods are unhealthy. TODO fix it right
				w.WriteHeader(http.StatusOK)
				return
			}
			deckURI := q.Get("deck")
			u, err := url.Parse(deckURI)
			if err != nil {

				log.Error(
					"Cannot parse deck uri",
					"err", err,
				)
				http.Error(w, "Cannot parse deck URI", http.StatusBadRequest)
				return
			}
			switch u.Host {
			//https://tappedout.net/mtg-decks/22-01-20-kess-storm/
			case "tappedout.net":
				deckURI = fmt.Sprintf("%s?fmt=txt", deckURI)
				log.Debug(
					"tappedout",
					"deckUri", deckURI,
				)
				var res *http.Response
				err := retry.Do(
					func() error {
						var err error
						res, err = http.DefaultClient.Get(deckURI)
						if err != nil {
							return err
						}
						return nil
					})
				if err != nil {
					log.Error(
						"cannot get deckbox deck",
						"err", err,
						"uri", deckURI,
					)
					http.Error(w, "Cannot get deckbox deck deck URI", http.StatusServiceUnavailable)
					return
				}
				if res.StatusCode != 200 {
					log.Error(
						"Unexpected response status",
						"status", res.Status,
					)
					http.Error(w, "Unexpected status code from deckbox", http.StatusBadGateway)
					return

				}
				r = res.Body
				break

			// https://deckbox.org/sets/2649137
			case "deckbox.org":
				deckURI = fmt.Sprintf("%s/export", deckURI)
				var res *http.Response
				err := retry.Do(
					func() error {
						var err error
						res, err = http.DefaultClient.Get(deckURI)
						if err != nil {
							return err
						}
						return nil
					})
				if err != nil {
					log.Error(
						"cannot get deckbox deck",
						"err", err,
						"uri", deckURI,
					)
					http.Error(w, "Cannot get deckbox deck deck URI", http.StatusServiceUnavailable)
					return
				}
				if res.StatusCode != 200 {
					log.Error(
						"Unexpected response status",
						"status", res.Status,
					)
					http.Error(w, "Unexpected status code from deckbox", http.StatusBadGateway)
					return

				}

				r, err = mtgfail.Normalize(res.Body, log)
				if err != nil {
					log.Error(
						"Unexpected format for deck status",
						"err", err,
						"url", deckURI,
					)
					http.Error(w, "Unexpected status code from deckbox", http.StatusBadGateway)
					return
				}
				break

			default:
				http.Error(w, "Unsupported deck host URI", http.StatusUnprocessableEntity)
				return
			}

		}

		deckList := make(map[string]int)

		switch req.Header.Get("Content-Type") {
		case "application/json":
			//TODO
		case "application/x-www-form-urlencoded":
			//TODO
		case "text/plain":
			r = req.Body
			fallthrough
		default:
			deckList, err = mtgfail.ReadCardList(r, log)
			if err != nil {
				log.Error(
					"Can't read cardfile",
					"err", err,
				)
				http.Error(w, "Can't read card list", http.StatusBadRequest)
				return
			}

		}

		deck, err := tts.BuildDeck(cache, deckList, log)
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
		w.Header().Add("Access-Control-Allow-Origin", "https:/api.mtg.fail")

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
