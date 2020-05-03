package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/avast/retry-go"
	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
)

var (
	def  string
	bulk string
)

func main() {

	flag.StringVar(&def, "file", "./deck.txt", "The fully qualified name of the deck definition")
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
	deckList, err := mtgfail.ReadCardList(def, log)
	if err != nil {
		log.Error(
			"Can't read cardfile",
			"err", err,
		)
		return
	}
	log.Info(
		"Read deck list",
		"entries", len(deckList),
	)

	deck := make(map[*mtgfail.Entry]int)

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
			err = retry.Do(
				func() error {
					var err error
					res, err = http.DefaultClient.Get(uri)
					if err != nil {
						log.Error(
							"cannot send request to scryfall",
							"err", err,
							"uri", uri,
						)
						return err
					}
					return nil
				})
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

		deck[entry] = count
		log.Info(
			"getting info for card name",
			"count", count,
			"name", name,
		)

	}

	d := buildStacks(log, deck)

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {

		log.Error(
			"Can't marshal deckfile",
			"err", err,
		)
		return

	}
	f, err := os.Create("out.json")
	if err != nil {

		log.Error(
			"Can't create file for deck",
			"err", err,
		)
		return

	}
	_, err = f.Write(b)
	if err != nil {
		log.Error(
			"Can't write deckfile",
			"err", err,
		)
		return

	}
	defer f.Close()

}

type AutoComplete struct {
	Object      string   `json:"object"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}

func buildStacks(log log15.Logger, stacks ...map[*mtgfail.Entry]int) *mtgfail.DeckFile {

	var (
		state []mtgfail.ObjectState
	)
	var deckNumber int
	cardNumber := 1
	for _, names := range stacks {

		log.Info(
			"building stack",
			"unique cards", len(names),
			"list", fmt.Sprintf("%+v", names),
		)
		var (
			ids  []int
			deck map[int]mtgfail.Card = make(map[int]mtgfail.Card, len(names))
			obs  []mtgfail.ContainedObject
		)
		for entry, count := range names {
			if entry == nil {

				log.Warn(
					"nil entry while building stack",
				)
				continue
			}

			for ; count > 0; count-- {
				log.Debug(
					"building card object card",
					"count", count,
					"name", entry.Name,
					"cardnumber", cardNumber,
				)

				cardTx := mtgfail.Transform{
					PosX:   0,
					PosY:   0,
					PosZ:   0,
					RotX:   0,
					RotY:   180,
					RotZ:   180,
					ScaleX: 1,
					ScaleY: 1,
					ScaleZ: 1,
				}

				id := (cardNumber) * 100
				ids = append(ids, id)
				ob := mtgfail.ContainedObject{
					CardID:    id,
					Name:      "Card",
					Nickname:  entry.Name,
					Transform: cardTx,
				}
				obs = append(obs, ob)
				card := mtgfail.Card{
					FaceURL:      entry.ImageUris.Normal,
					BackURL:      "https://www.frogtown.me/images/gatherer/CardBack.jpg",
					NumHeight:    1,
					NumWidth:     1,
					BackIsHidden: true,
				}

				deck[cardNumber] = card

				cardNumber++
			}
			deckNumber++
		}
		var z int
		if deckNumber == 0 {
			z = 0 // face down
		} else {

			z = 180 // face up
		}

		stackTx := mtgfail.Transform{
			PosX:   deckNumber + 2,
			PosY:   1,
			PosZ:   0,
			RotX:   0,
			RotY:   180,
			RotZ:   z,
			ScaleX: 1,
			ScaleY: 1,
			ScaleZ: 1,
		}
		state = append(state, mtgfail.ObjectState{
			Name:             "DeckCustom",
			ContainedObjects: obs,
			CustomDeck:       deck,
			DeckIDs:          ids,
			Transform:        stackTx,
		})
	}
	return &mtgfail.DeckFile{
		ObjectStates: state,
	}
}
