package mtgfail

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/inconshreveable/log15"
)

// ExampleDeck ...
const ExampleDeck = "examples/deck.txt"

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
