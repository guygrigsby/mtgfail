package tabletopsimulator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/getlantern/deepcopy"
	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
)

// BuildDeck ...
func BuildDeck(ctx context.Context, bulk mtgfail.CardStore, deckList map[string]int, log log15.Logger) (*DeckFile, error) {
	deck := make(map[*mtgfail.Entry]int)
	tokenDeck := make(map[*mtgfail.Entry]int)

	for name, count := range deckList {
		entry, err := bulk.Get(name, log)
		if err != nil {
			log.Error(
				"failed to conact store to get card",
				"err", err,
			)
			return nil, err
		}
		if entry == nil {
			log.Warn(
				"cache miss. Calling scryfall for autocomplete",
				"name", name,
				"count", count,
				"err", err,
			)

			entry, err = findAlternateNamedCards(name, bulk, log)
			if err != nil {
				msg := "failed to contact scryfall and get corrected card"
				log.Error(
					msg,
					"err", err,
				)
				return nil, err
			}

		}
		if isDoubleSided(entry) {
			token, err := createTokenEntry(*entry, log)
			if err != nil {
				return nil, err
			}
			tokenDeck[entry] = count
			entry = token

		}

		deck[entry] = count

	}
	if len(tokenDeck) > 0 {
		return BuildStacks(log, deck, tokenDeck)
	}
	return BuildStacks(log, deck)

}

func isDoubleSided(entry *mtgfail.Entry) bool {
	return strings.Contains(entry.Name, "//")
}

func createTokenEntry(entry mtgfail.Entry, log log15.Logger) (*mtgfail.Entry, error) {
	front := strings.Split(entry.CardFaces[0].ImageUris.Large, "?")[0]
	back := strings.Split(entry.CardFaces[1].ImageUris.Large, "?")[0]
	var token mtgfail.Entry
	err := deepcopy.Copy(token, entry)
	if err != nil {
		log.Error(
			"Could not create copy of double sided card",
			"cardname", entry.Name,
			"err", err,
		)
		return nil, err
	}

	log.Info(
		"Double sided card",
		"name", entry.Name,
		"front", front,
		"back", back,
		"token", token,
	)
	return &token, nil
}

func findAlternateNamedCards(name string, bulk mtgfail.CardStore, log log15.Logger) (*mtgfail.Entry, error) {
	escName := url.QueryEscape(name)
	uri := fmt.Sprintf("https://api.scryfall.com/cards/autocomplete?q=%s", escName)
	var entry *mtgfail.Entry

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
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Error(
			"Unexpected response status",
			"status", res.Status,
		)

		return nil, err

	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(
			"cannot read scryfall response body",
			"err", err,
			"uri", uri,
		)
		return nil, err
	}
	var autoComplete AutoComplete

	err = json.Unmarshal(b, &autoComplete)
	if err != nil {
		log.Error(
			"Cannot unmarshal scryfal res",
			"err", err,
		)
		return nil, err
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

	// set it in the local data
	err = bulk.Put(name, entry, log)
	if err != nil {
		log.Error(
			"cannot put correct card in store. Continuing without saving",
			"err", err,
		)
		// we can continue
	}
	return entry, nil

}

type AutoComplete struct {
	Object      string   `json:"object"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}

func BuildStacks(log log15.Logger, stacks ...map[*mtgfail.Entry]int) (*DeckFile, error) {

	var (
		state []ObjectState
	)
	var deckNumber int
	cardNumber := 1
	for _, stack := range stacks {

		log.Info(
			"building stack",
			"unique cards", len(stack),
		)
		var (
			deck             = make(map[int]Card, len(stack))
			ids              []int
			containedObjects []ContainedObject
		)
		var cardCount int
		for entry, v := range stack {
			if v == 0 {
				log.Warn(
					"Encountered card with 0 occurrences in deck count card. Assuming 1.",
					"cardname", entry.Name,
				)
				v = 1
			}
			cardCount += v
			if entry == nil {

				log.Warn(
					"nil entry while building stack",
				)
				continue
			}
			// Card multiples
			for count := v; count > 0; count-- {

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
				containedObjects = append(containedObjects, ob)
				var img string

				if isDoubleSided(entry) {
					img = entry.CardFaces[0].ImageUris.Large
				} else {
					img = entry.ImageUris.Large
				}

				card := Card{
					FaceURL:      img,
					BackURL:      "https://firebasestorage.googleapis.com/v0/b/marketplace-c87d0.appspot.com/o/card_back.jpg?alt=media",
					NumHeight:    1,
					NumWidth:     1,
					BackIsHidden: true,
				}
				log.Info(
					"Card Created",
					"card", card,
				)

				deck[cardNumber] = card

				cardNumber++
			}

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
			ContainedObjects: containedObjects,
			CustomDeck:       deck,
			DeckIDs:          ids,
			Transform:        stackTx,
		})

		deckNumber++
	}
	return &DeckFile{
		ObjectStates: state,
	}, nil
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
	ContainedObjects []ContainedObject `json:"ContainedObjects,omitempty"`
	CustomDeck       map[int]Card      `json:"CustomDeck"`
	DeckIDs          []int             `json:"DeckIDs,omitempty"`
	Transform        Transform         `json:"Transform"`
	CardID           int               `json:"CardID,omitempty"`
	Nickname         string            `json:"Nickname,omitempty"`
}
