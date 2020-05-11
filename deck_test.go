package mtgfail_test

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/guygrigsby/mtgfail"
	tts "github.com/guygrigsby/mtgfail/tabletopsimulator"
	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
)

var (
	text string
	bulk string
)

func TestDesk(t *testing.T) {

	flag.StringVar(&text, "file", mtgfail.ExampleDeck, "The fully qualified name of the deck definition")
	flag.StringVar(&bulk, "bulk", "./scryfall-default-cards.json", "The bulk data download")

	log := log15.New()

	bulk, err := mtgfail.ReadBulk(bulk, log)
	assert.NoError(t, err)

	r, err := os.Open(text)
	assert.NoError(t, err)

	deckList, err := mtgfail.ReadCardList(r, log)
	assert.NoError(t, err)

	deck, err := tts.BuildDeck(context.Background(), bulk, deckList, log)
	assert.NoError(t, err)

	dumpToFile(deck, "golden-deck.json")
	if true {
		return
	}

	f, err := os.Open("test-out.json")
	assert.NoError(t, err)
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	var golden tts.DeckFile

	err = json.Unmarshal(b, &golden)
	assert.NoError(t, err)

	assert.Equal(t, golden.ObjectStates[0].DeckIDs, deck.ObjectStates[0].DeckIDs, "Deck does not match golden example")

	for _, expected := range golden.ObjectStates[0].ContainedObjects {
		found := false
		var got tts.ContainedObject

		for _, got = range deck.ObjectStates[0].ContainedObjects {
			if expected.Nickname == got.Nickname {
				found = true
				expectedCard := golden.ObjectStates[0].CustomDeck[expected.CardID]
				gotCard := deck.ObjectStates[0].CustomDeck[got.CardID]
				assert.Equal(t, expectedCard, gotCard, "Card mismatch")
				continue
			}

		}
		assert.True(t, found, "Did not find required Card in contained objects", got.Nickname)

	}

	defer f.Close()

}

func dumpToFile(i interface{}, name string) error {

	b, err := json.Marshal(i)
	if err != nil {

		return err
	}

	f, err := os.Create(name)
	if err != nil {

		return err
	}
	_, err = f.Write(b)
	if err != nil {

		return err
	}

	return f.Close()
}
