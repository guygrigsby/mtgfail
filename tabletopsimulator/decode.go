package tabletopsimulator

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
)

// Decode ...
func Decode(ttsFormat io.Reader, log log15.Logger) ([]*mtgfail.Entry, error) {
	b, err := ioutil.ReadAll(ttsFormat)
	if err != nil {

		log.Error(
			"Cannot read deckfile",
			"err", err,
		)
		return nil, err
	}
	var deckfile DeckFile
	err = json.Unmarshal(b, &deckfile)
	if err != nil {
		log.Error(
			"Cannt unmarshal tts deckfile",
			"err", err,
		)
		return nil, err
	}

	var deck []*mtgfail.Entry

	for _, entry := range deckfile.ObjectStates {

		for i, obj := range entry.ContainedObjects {
			deck = append(deck, &mtgfail.Entry{
				Name: obj.Nickname,
				ImageUris: mtgfail.ImageUris{
					Normal: entry.CustomDeck[i+1].FaceURL,
					Small:  entry.CustomDeck[i+1].FaceURL,
					Large:  entry.CustomDeck[i+1].FaceURL,
					Png:    entry.CustomDeck[i+1].FaceURL,
				},
			})
		}

	}

	return deck, nil
}
