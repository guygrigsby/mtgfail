package mtgfail

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/inconshreveable/log15"
)

// ExampleDeck ...
const ExampleDeck = "examples/deck.txt"

// ReadBulk ...
func ReadBulk(file string, log log15.Logger) (CardStore, error) {
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
	var bulk = store(make(map[string]*Entry))
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

type store map[string]*Entry

func (s store) GetMany(names []string, log log15.Logger) ([]*Entry, error) {
	var matches []*Entry
	for _, name := range names {
		m, ok := s[name]
		if ok {
			matches = append(matches, m)
		}
	}
	return matches, nil
}
func (s store) Get(name string, log log15.Logger) (*Entry, error) {
	return s[name], nil
}
func (s store) Put(name string, e *Entry, log log15.Logger) error {
	s[name] = e
	return nil
}
func (s store) PutMany(entries map[string]*Entry, log log15.Logger) error {
	for k, v := range entries {
		s[k] = v
	}
	return nil
}

// ConvertToPairText ...
func ConvertToPairText(deck *Deck) (map[string]int, error) {
	cards := make(map[string]int)
	if len(deck.Cards) == 0 {
		return nil, fmt.Errorf("Zero length deck %+v", deck)
	}
	for _, card := range deck.Cards {
		count := cards[card.Name]
		count++
		cards[card.Name] = count
	}
	return cards, nil
}

func ConvertEntriesToPairText(cards []*Entry) (map[string]int, error) {
	pairs := make(map[string]int)
	if len(cards) == 0 {
		return nil, fmt.Errorf("Zero length deck %+v", cards)
	}
	for _, card := range cards {
		count := pairs[card.Name]
		count++
		pairs[card.Name] = count
	}
	return pairs, nil
}

// ReadCardList ...
func ReadCardList(r io.ReadCloser, log log15.Logger) (map[string]int, error) {

	cards := make(map[string]int)
	scanner := bufio.NewScanner(r)
	defer r.Close()

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		lineScanner := bufio.NewScanner(strings.NewReader(line))
		lineScanner.Split(bufio.ScanWords)
		ok := lineScanner.Scan()
		if !ok {
			continue
		}
		str := lineScanner.Text()
		count, err := strconv.Atoi(str)
		if err != nil {
			// Assume one if there is no number
			log.Debug(
				"Error converting count",
				"count", count,
				"str", str,
			)
			count = 1
		}
		sb := strings.Builder{}

		for lineScanner.Scan() {
			txt := lineScanner.Text()
			if txt == "X" || txt == "x" {
				log.Debug(
					"throwing away x",
					"text", txt,
				)
				continue
			}
			sb.WriteString(strings.ToLower(txt))
			sb.WriteString(" ")

		}
		name := strings.TrimSpace(sb.String())
		if count == 0 {
			// We'll assume 1 and fix this, but there is probably a root cause that we are missing.
			// And we can't be the source of truth, for now.
			count = 1
			log.Warn(
				"Found zero count card, changing to 1",
				"cardname", name,
			)

		}
		log.Debug(
			"adding",
			"name", name,
			"count", count,
		)
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
