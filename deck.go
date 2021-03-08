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

func NormalizeCardName(name string) string {
	//if strings.Contains(name, "//") {
	//	name = strings.ReplaceAll(name, "//", "")
	//	re := regexp.MustCompile(`//.*`)
	//	// Strip everything after the double slash
	//	// Scrycall has the // and that's where we get our card data
	//	return string(re.ReplaceAll([]byte(name), []byte{}))

	//}
	return name
}

// BuildDeck ...
func BuildDeck(ctx context.Context, bulk CardStore, deckList map[string]int, log log15.Logger) (*Deck, error) {
	var (
		deck = Deck{
			Cards: nil,
		}
	)
	for name, count := range deckList {
		entry, err := bulk.Get(Key(name), log)
		if entry == nil || err != nil {
			log.Warn(
				"cache miss. Calling scryfall for autocomplete",
				"name", name,
				"count", count,
				"err", err,
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
			entry, err = bulk.Get(correctName, log)
			if err != nil {
				log.Error(
					"cannot access store",
					"err", err,
				)
				return nil, err
			}
			err = bulk.Put(name, entry, log)
			if err != nil {
				log.Error(
					"cannot put store",
					"err", err,
				)
				return nil, err
			}
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
