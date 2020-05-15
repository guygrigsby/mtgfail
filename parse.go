package mtgfail

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/inconshreveable/log15"
)

const ExampleDeck = "examples/deck.txt"

func ReadBulk(file string, log log15.Logger) (Bulk, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Error(
			"Can't open file",
			"err", err,
		)
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(
			"Can't read file",
			"err", err,
		)
		return nil, err
	}

	var cards []*Entry
	err = json.Unmarshal(b, &cards)
	if err != nil {
		log.Error(
			"Can't unmarshal data",
			"err", err,
		)
		return nil, err
	}
	var bulk = make(map[string]*Entry)
	for i, card := range cards {
		if card == nil {
			log.Warn(
				"nil entry skipping",
				"index", i,
			)
			continue
		}
		//TODO it's gross, but scryfall adds the time of download as a param at the end and tts no likey
		card.ImageUris.Small = strings.Split(card.ImageUris.Small, "?")[0]
		card.ImageUris.Normal = strings.Split(card.ImageUris.Normal, "?")[0]
		card.ImageUris.Large = strings.Split(card.ImageUris.Large, "?")[0]
		card.ImageUris.Png = strings.Split(card.ImageUris.Png, "?")[0]
		bulk[card.Name] = card

	}

	return bulk, nil
}

// ReadCardList
func ReadCardList(r io.ReadCloser, log log15.Logger) (map[string]int, error) {

	//b, _ := ioutil.ReadAll(r)
	log.Debug(
		"scanning ",
	//	"content", string(b),
	)
	cards := make(map[string]int)
	scanner := bufio.NewScanner(r)
	defer r.Close()

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		log.Debug(
			"scanning line",
			"val", line,
		)

		lineScanner := bufio.NewScanner(strings.NewReader(line))
		lineScanner.Split(bufio.ScanWords)
		ok := lineScanner.Scan()
		if !ok {
			log.Info(
				"Cannot scan count from start of line",
				"line", line,
			)
			continue
		}
		str := lineScanner.Text()
		count, err := strconv.Atoi(str)
		if err != nil {
			log.Error(
				"Invalid file format. Cannot parse card count.",
				"err", err,
				"val", str,
			)
			return nil, err
		}
		sb := strings.Builder{}
		for lineScanner.Scan() {
			txt := lineScanner.Text()
			log.Debug(
				"scanning word token name",
				"val", txt,
			)
			sb.WriteString(txt)
			sb.WriteString(" ")

		}
		name := strings.TrimSpace(sb.String())
		cards[name] = count

	}
	if err := scanner.Err(); err != nil {
		log.Error(
			"Scanner error",
			"err", err,
		)
	}
	return cards, nil
}
