package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/guygrigsby/mtgfail"
)

func main() {

	deck := []string{
		"testcarda",
		"cardw",
		"stestcarda",
	}
	tokens := []string{
		"token1",
		"token2",
	}

	d := buildStacks(deck, tokens)

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)

	}
	f, err := os.Create("out.json")
	if err != nil {
		panic(err)

	}
	_, err = f.Write(b)
	if err != nil {
		panic(err)

	}
	defer f.Close()

}

func buildStacks(stacks ...[]string) *mtgfail.DeckFile {
	var (
		state []mtgfail.ObjectState
	)

	for i, names := range stacks {
		var (
			ids  []int
			deck map[int]mtgfail.Card = make(map[int]mtgfail.Card, len(names))
			obs  []mtgfail.ContainedObject
		)
		for i, name := range names {
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

			id := (i + 1) * 100
			ids = append(ids, id)
			ob := mtgfail.ContainedObject{
				CardID:    id,
				Name:      "Card",
				Nickname:  name,
				Transform: cardTx,
			}
			obs = append(obs, ob)
			card := mtgfail.Card{
				FaceURL:      "https://www.frogtown.me/Images/1b6772e4-4ff9-4090-a8dc-df9e8f568fc0.jpg",
				BackURL:      "https://www.frogtown.me/images/gatherer/CardBack.jpg",
				NumHeight:    1,
				NumWidth:     1,
				BackIsHidden: true,
			}

			deck[i+1] = card

		}
		var z int
		if i == 0 {
			z = 180 // face down
		} else {

			z = 0 // face up
		}

		stackTx := mtgfail.Transform{
			PosX:   i + 2,
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

func thing() {
	ctx := context.Background()
	client, err := scryfall.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	sco := scryfall.SearchCardsOptions{
		Unique:        scryfall.UniqueModePrints,
		Order:         scryfall.OrderSet,
		Dir:           scryfall.DirDesc,
		IncludeExtras: true,
	}
	result, err := client.SearchCards(ctx, "storm cro", sco)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", result.Cards[0].Colors)
}
