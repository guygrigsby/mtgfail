package mtgfail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go"
	"github.com/inconshreveable/log15"
)

type AutoComplete struct {
	Object      string   `json:"object"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}
type CardShort struct {
	Name   string
	Cost   string
	Cmc    float64
	Image  string
	Rarity string
	Set    string
	Colors []string
	Text   string
}

// BuildDeck ...
func BuildDeck(ctx context.Context, bulk Bulk, deckList map[string]int, log log15.Logger) (*Deck, error) {
	var (
		deck = Deck{
			Cards: nil,
		}
	)

	for name, count := range deckList {
		entry := bulk[name]
		if entry == nil {
			log.Warn(
				"cache miss. Calling scryfall for autocomplete",
				"name", name,
				"count", count,
			)
			escName := url.QueryEscape(name)
			uri := fmt.Sprintf("https://api.scryfall.com/cards/autocomplete?q=%s", escName)

			var res *http.Response
			err := retry.Do(
				func() error {
					var err error
					c := http.Client{Timeout: 5 * time.Second}
					res, err = c.Get(uri)
					if err != nil {
						return err
					}

					return nil
				})

			if err != nil {
				log.Error(
					"cannot send request to scryfall",
					"err", err,
					"uri", uri,
				)
				continue
			}

			if res.StatusCode != 200 {
				log.Error(
					"Unexpected response status",
					"status", res.Status,
				)
				continue

			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Error(
					"cannot read scryfall response body",
					"err", err,
					"uri", uri,
				)
				continue
			}

			var autoComplete AutoComplete
			err = json.Unmarshal(b, &autoComplete)
			if err != nil {
				log.Error(
					"Cannot unmarshal scryfal res",
					"err", err,
				)
				continue
			}

			correctName := autoComplete.Data[0]
			entry = bulk[correctName]
			bulk[name] = entry
			log.Debug(
				"Scryfall autocomplete success",
				"original", name,
				"corrected", correctName,
			)

		}

		card := &CardShort{
			Name:   entry.Name,
			Cost:   entry.ManaCost,
			Cmc:    entry.Cmc,
			Image:  entry.ImageUris.Large,
			Rarity: entry.Rarity,
			Set:    entry.Set,
			Colors: entry.Colors,
			Text:   entry.OracleText,
		}

		for i := 0; i < count; i++ {
			deck.Cards = append(deck.Cards, card)
		}

	}
	return &deck, nil

}
