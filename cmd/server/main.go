package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go"
	"github.com/guygrigsby/mtgfail"
	tts "github.com/guygrigsby/mtgfail/pkg/tabletopsimulator"
	"github.com/inconshreveable/log15"
)

// BuildDeck ...
func BuildDeck(cache mtgfail.Bulk, log log15.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Debug(
			"Request",
			"method", req.Method,
			"caller", req.RemoteAddr,
			"headers", fmt.Sprintf("%+v", req.Header),
		)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Cache-Control", "no-cache")

		if req.Method == "OPTIONS" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("OKOK")
			return
		}
		var (
			err     error
			content io.ReadCloser = req.Body
		)
		ctx, cancel := context.WithTimeout(req.Context(), time.Second*60)
		defer cancel()
		if req.Method == http.MethodGet {
			q := req.URL.Query()
			if len(q) == 0 {
				// GCP load balancer health checks are garbage. Somehow, they always end up at '/'
				// This was I don' spend hours softing out why my pods are unhealthy. TODO fix it right
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("No params")
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
						c := http.Client{
							Timeout: 5 * time.Second,
						}
						res, err = c.Get(deckURI)
						if err != nil {
							return err
						}
						return nil
					},
					retry.Attempts(3),
				)
				if err != nil {
					log.Error(
						"cannot get tappedout deck",
						"err", err,
						"uri", deckURI,
					)
					http.Error(w, "Cannot get tappedout deck URI", http.StatusServiceUnavailable)
					return
				}
				if res.StatusCode != 200 {
					log.Error(
						"Unexpected response status",
						"status", res.Status,
					)
					http.Error(w, "Unexpected status code from tappedout", http.StatusBadGateway)
					return

				}
				content = res.Body
				break

			// https://deckbox.org/sets/2649137
			case "deckbox.org":
				deckURI = fmt.Sprintf("%s/export", deckURI)
				log.Debug(
					"deckbox",
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

				content, err = mtgfail.Normalize(res.Body, log)
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
			msg := "application/json not supported"
			log.Error(
				msg,
			)
			http.Error(w, "Empty deck list", http.StatusUnsupportedMediaType)
			return
		case "application/x-www-form-urlencoded":
			//TODO
			msg := "application/x-www-form-url-encoded  not supported"
			log.Error(
				msg,
			)
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		case "text/plain":
			b, err := ioutil.ReadAll(content)
			if err != nil {
				log.Error(
					"error reading body",
					"err", err,
				)
				http.Error(w, "Can't read body", http.StatusInternalServerError)
				return

			}
			req.Body.Close()
			deckList, err = mtgfail.ReadCardList(
				ioutil.NopCloser(
					bytes.NewBuffer(b)), log)
			if err != nil {
				log.Error(
					"Can't read cardfile",
					"err", err,
				)
				http.Error(w, "Can't read card list", http.StatusBadRequest)
				return
			}
			if len(deckList) == 0 {
				log.Error(
					"Zero length decklist",
					"content", string(b),
				)
				http.Error(w, "Empty deck list", http.StatusInternalServerError)
				return

			}
		default:

			msg := fmt.Sprintf("Unrecognized content type %s", req.Header.Get("Content-Type"))
			log.Error(
				msg,
			)
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return

		}

		deck, err := tts.BuildDeck(ctx, cache, deckList, log)
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
