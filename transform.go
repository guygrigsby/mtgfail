package mtgfail

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/inconshreveable/log15"
	"golang.org/x/net/html"
)

/*
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <meta name="description" content="Card set print view" />
    <title>Deckbox.org card set print view</title>
</head>
<body>
        1 Brainstorm<br/>1 Thrasios, Triton Hero<br/>

    <p><strong>Sideboard:</strong></p>


</body>
</html>
*/
type DeckSite int

const (
	Scryfall DeckSite = iota
	DeckBox
	TappedOut
)

// Normalize ...
func Normalize(source DeckSite, r io.ReadCloser, log log15.Logger) (map[string]int, error) {
	var (
		deck map[string]int
		err  error
	)
	switch source {
	case DeckBox:
		log.Info("deckbox deck")
		deck, err = normalizeDeckbox(r, log)
	case TappedOut:
		log.Info("tappedout deck")
		deck, err = normalizeTappedOut(r, log)
	case Scryfall:
		log.Debug("scryfall deck")
		deck, err = normalizeScryfall(r, log)
	}
	log.Debug(
		"normalized",
		"content", deck,
	)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func normalizeScryfall(r io.ReadCloser, log log15.Logger) (map[string]int, error) {
	return ReadCardList(r, log)
}

func normalizeTappedOut(r io.ReadCloser, log log15.Logger) (map[string]int, error) {
	return ReadCardList(r, log)

}
func normalizeDeckbox(r io.ReadCloser, log log15.Logger) (map[string]int, error) {

	z := html.NewTokenizer(r)

	w := strings.Builder{}
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		t := z.Token()
		if tt == html.StartTagToken && t.Data == "body" {
			log.Debug("body", "data", t.Data)

			for tt = z.Next(); tt != html.ErrorToken; tt = z.Next() {
				t := z.Token()

				if tt == html.EndTagToken && t.Data == "body" ||
					tt == html.StartTagToken && t.Data == "p" {
					return ReadCardList(ioutil.NopCloser(strings.NewReader(w.String())), log)
				}

				if tt == html.SelfClosingTagToken && t.Data == "br" {
					_, err := w.WriteRune('\n')
					if err != nil {
						log.Error(
							"Cannot write newline",
							"err", err,
						)
						return nil, err
					}

				} else if tt == html.EndTagToken {
					log.Debug("html end token", "data", t.Data)
					continue
				} else {
					_, err := w.WriteString(t.Data)
					if err != nil {
						log.Error(
							"Cannot write data",
							"err", err,
						)
						return nil, err

					}
				}
			}
		}
	}
	return ReadCardList(ioutil.NopCloser(strings.NewReader(w.String())), log)

}
