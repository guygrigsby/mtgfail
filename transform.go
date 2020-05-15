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

// Normalize ...
func Normalize(r io.ReadCloser, log log15.Logger) (io.ReadCloser, error) {

	z := html.NewTokenizer(r)

	w := strings.Builder{}
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		t := z.Token()
		if tt == html.StartTagToken && t.Data == "body" {
			log.Debug("body", "data", t.Data)

			for tt = z.Next(); tt != html.ErrorToken; tt = z.Next() {
				t := z.Token()
				log.Info("token", "data", t.Data)

				if tt == html.EndTagToken && t.Data == "body" ||
					tt == html.StartTagToken && t.Data == "p" {
					return ioutil.NopCloser(strings.NewReader(w.String())), nil
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
					log.Info("writing token", "data", t.Data)
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
	log.Debug(
		"normalized",
		"content", w.String(),
	)
	return ioutil.NopCloser(strings.NewReader(w.String())), nil
}
