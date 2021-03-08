package deck

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
	"golang.org/x/sync/errgroup"

	tts "github.com/guygrigsby/mtgfail/tabletopsimulator"
)

type Format int

const (
	MaxAllowedContentLength = 50000

	TableTopSimulator Format = iota
	ScryfallEntry
	All
)

func FetchDeck(deckURI string, log log15.Logger) (io.ReadCloser, error, int) {
	var (
		err     error
		content io.ReadCloser
	)

	u, err := url.Parse(deckURI)
	if err != nil {

		log.Error(
			"Cannot parse deck uri",
			"err", err,
		)
		return nil, err, http.StatusBadRequest
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
			return nil, fmt.Errorf("Cannot get tappedout deck: %w", err), http.StatusServiceUnavailable
		}
		if res.StatusCode != 200 {
			log.Error(
				"Unexpected response status",
				"status", res.Status,
			)
			return nil, fmt.Errorf("Unexpected status code from tappedout"), http.StatusBadRequest

		}
		content = res.Body

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
			return nil, err, http.StatusServiceUnavailable
		}
		if res.StatusCode != 200 {
			log.Error(
				"Unexpected response status",
				"status", res.Status,
			)
			return nil, fmt.Errorf("unexpected status %v", res.StatusCode), http.StatusBadGateway

		}

		content, err = mtgfail.Normalize(res.Body, log)
		if err != nil {
			log.Error(
				"Unexpected format for deck status",
				"err", err,
				"url", deckURI,
			)
			return nil, err, http.StatusBadGateway
		}

	default:
		log.Debug(
			"Unexpected deck Host",
			"url", deckURI,
			"Host", u.Host,
		)

		return nil, fmt.Errorf("Unknown Host"), http.StatusUnprocessableEntity
	}
	return content, nil, 200
}

// BuildTTSDeck ...
func BuildTTSDeck(cache mtgfail.CardStore, log log15.Logger) http.HandlerFunc {
	return BuildDeck(TableTopSimulator, cache, log)
}

// BuildInternalDeck ...
func BuildInternalDeck(cache mtgfail.CardStore, log log15.Logger) http.HandlerFunc {
	return BuildDeck(ScryfallEntry, cache, log)
}

// BuildAll ...
func BuildAll(cache mtgfail.CardStore, log log15.Logger) http.HandlerFunc {
	return BuildDeck(All, cache, log)
}

// BuildDeck ...
func BuildDeck(f Format, cache mtgfail.CardStore, log log15.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Debug(
			"Request",
			"method", req.Method,
			"caller", req.RemoteAddr,
			"params", fmt.Sprintf("%+v", req.URL.Query()),
			"headers", fmt.Sprintf("%+v", req.Header),
			"req", fmt.Sprintf("%+v", req),
		)

		if req.Header.Get("Expect") == "100-continue" {
			l, err := strconv.Atoi(req.Header.Get("Content-Length"))
			if err != nil {
				log.Error(
					"cannot parse content length for 100-continue",
					"err", err,
				)
				http.Error(w, err.Error(), http.StatusExpectationFailed)
			}
			if l > MaxAllowedContentLength {
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			w.WriteHeader(http.StatusContinue)
			return
		}
		w.Header().Add("Cache-Control", "no-cache")

		var (
			err     error
			content io.ReadCloser = req.Body
		)
		ctx, cancel := context.WithTimeout(req.Context(), time.Second*60)
		defer cancel()
		switch {
		case req.Method == http.MethodGet:
			var status int
			q, err := url.ParseQuery(req.URL.RawQuery)
			if err != nil {
				log.Error(
					"unable to parse query params",
					"err", err,
				)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if len(q) == 0 {
				// GCP load balancer health checks are garbage.
				// Somehow, they always end up at '/'
				return
			}

			content, err, status = FetchDeck(q.Get("deck"), log)
			if status > 299 {
				log.Error(
					"unexpected return status",
					"status", status,
				)
				http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
				return
			}

			if err != nil {
				log.Error(
					"failed to fetch deck",
					"err", err,
				)
				http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
				return

			}

		case req.Method == http.MethodPost:
			// leave the body alone!
			content = req.Body

		}

		if content == nil {
			log.Error(
				"failed to get content",
			)
			http.Error(w, "No content", http.StatusInternalServerError)
			return
		}

		deckList := make(map[string]int)

		switch req.Header.Get("Content-Type") {
		case "application/json":
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
			var deck []*mtgfail.Entry
			err = json.Unmarshal(b, &deck)
			if err != nil {
				msg := "cannot unmarshal JSON deck"
				log.Error(
					msg,
				)

				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			log.Debug(
				"json body",
				"body", string(b),
				"marshalled", fmt.Sprintf("%+v", deck),
			)

			deckList, err = mtgfail.ConvertEntriesToPairText(deck)
			if err != nil {
				msg := "empty deck"
				log.Error(
					msg,
				)

				http.Error(w, msg, http.StatusBadRequest)
				return
			}

		case "text/plain", "":
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
		var deck interface{}

		switch f {
		case TableTopSimulator:
			deck, err = tts.BuildDeck(ctx, cache, deckList, log)
		case ScryfallEntry:
			deck, err = mtgfail.BuildDeck(ctx, cache, deckList, log)
		case All:
			ret := &DualDeck{}
			g, ctx := errgroup.WithContext(ctx)
			g.Go(func() error {
				ttsDeck, err := tts.BuildDeck(ctx, cache, deckList, log)
				if err != nil {
					return err
				}
				ret.TTS = ttsDeck
				return nil
			})
			g.Go(func() error {
				internDeck, err := mtgfail.BuildDeck(ctx, cache, deckList, log)
				if err != nil {
					return err
				}
				ret.Intern = internDeck
				return nil
			})
			if err := g.Wait(); err != nil {
				log.Error(
					"Cannot build decks",
					"err", err,
				)
				return
			}
			b, err := json.Marshal(ret)
			if err != nil {

				log.Error(
					"Can't marshal deckfile",
					"err", err,
				)
				return

			}

			w.Header().Add("Content-Type", "application/json")

			_, err = w.Write(b)
			if err != nil {
				log.Error(
					"Can't write dual deckfile",
					"err", err,
				)
				return

			}

		}

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

type DualDeck struct {
	TTS    *tts.DeckFile `json:"tts,omitempty"`
	Intern *mtgfail.Deck `json:"internal,omitempty"`
}

func keySet(m map[string]int) []string {
	var s []string
	for k := range m {
		s = append(s, k)
	}
	return s
}
