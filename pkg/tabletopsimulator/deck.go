package tabletopsimulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/avast/retry-go"
	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
)

func BuildDeck(bulk mtgfail.Bulk, deckList map[string]int, log log15.Logger) (*DeckFile, error) {
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
			err := retry.Do(
				func() error {
					var err error
					res, err = http.DefaultClient.Get(uri)
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

		deck[entry] = count
		log.Info(
			"getting info for card name",
			"count", count,
			"name", name,
		)

	}

	return buildStacks(log, deck), nil

}

type AutoComplete struct {
	Object      string   `json:"object"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}

func buildStacks(log log15.Logger, stacks ...map[*mtgfail.Entry]int) *DeckFile {

	var (
		state []ObjectState
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
			deck map[int]Card = make(map[int]Card, len(names))
			obs  []ContainedObject
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

				cardTx := Transform{
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
				ob := ContainedObject{
					CardID:    id,
					Name:      "Card",
					Nickname:  entry.Name,
					Transform: cardTx,
				}
				obs = append(obs, ob)
				card := Card{
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

		stackTx := Transform{
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
		state = append(state, ObjectState{
			Name:             "DeckCustom",
			ContainedObjects: obs,
			CustomDeck:       deck,
			DeckIDs:          ids,
			Transform:        stackTx,
		})
	}
	return &DeckFile{
		ObjectStates: state,
	}
}

type DeckFile struct {
	ObjectStates []ObjectState `json:"ObjectStates"`
}

type Card struct {
	FaceURL      string `json:"FaceURL"`
	BackURL      string `json:"BackURL"`
	NumHeight    int    `json:"NumHeight"`
	NumWidth     int    `json:"NumWidth"`
	BackIsHidden bool   `json:"BackIsHidden"`
}
type Transform struct {
	PosX   int `json:"posX"`
	PosY   int `json:"posY"`
	PosZ   int `json:"posZ"`
	RotX   int `json:"rotX"`
	RotY   int `json:"rotY"`
	RotZ   int `json:"rotZ"`
	ScaleX int `json:"scaleX"`
	ScaleY int `json:"scaleY"`
	ScaleZ int `json:"scaleZ"`
}

type ContainedObject struct {
	CardID    int       `json:"CardID"`
	Name      string    `json:"Name"`
	Nickname  string    `json:"Nickname"`
	Transform Transform `json:"Transform"`
}

type ObjectState struct {
	Name             string            `json:"Name"`
	ContainedObjects []ContainedObject `json:"ContainedObjects"`
	CustomDeck       map[int]Card      `json:"CustomDeck"`
	DeckIDs          []int             `json:"DeckIDs"`
	Transform        Transform         `json:"Transform"`
}
