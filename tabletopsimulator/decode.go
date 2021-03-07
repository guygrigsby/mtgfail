package tabletopsimulator

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/guygrigsby/mtgfail"
	"github.com/inconshreveable/log15"
)

func Decode(ttsFormat io.Reader, log log15.Logger) ([]*mtgfail.Entry, error) {
	b, err := ioutil.ReadAll(ttsFormat)
	if err != nil {

		log.Error(
			"Cannot read deckfile",
			"err", err,
		)
		return nil, err
	}
	var deckfile []ObjectState
	err = json.Unmarshal(b, &deckfile)
	if err != nil {
		log.Error(
			"Cannt unmarshal tts deckfile",
			"err", err,
		)
		return nil, err
	}

	var deck []*mtgfail.Entry

	for i, entry := range deckfile {

		deck = append(deck, &mtgfail.Entry{
			Name: entry.Nickname,
			ImageUris: mtgfail.ImageUris{
				Normal: entry.CustomDeck[i+1].FaceURL,
			},
		})
	}

	return deck, nil
}
