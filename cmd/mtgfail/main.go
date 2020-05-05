package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/guygrigsby/mtgfail"
	tts "github.com/guygrigsby/mtgfail/pkg/tabletopsimulator"
	"github.com/inconshreveable/log15"
)

var (
	def  string
	bulk string
)

func main() {

	flag.StringVar(&def, "file", mtgfail.ExampleDeck, "The fully qualified name of the deck definition")
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
	r, err := os.Open(def)
	if err != nil {
		log.Error(
			"Can't open file",
			"err", err,
		)
		return
	}

	deckList, err := mtgfail.ReadCardList(r, log)
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

	deck, err := tts.BuildDeck(bulk, deckList, log)
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
