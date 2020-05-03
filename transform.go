package mtgfail

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/inconshreveable/log15"
	"golang.org/x/net/html"
)

func Normalize(r io.ReadCloser, log log15.Logger) (io.ReadCloser, error) {

	z := html.NewTokenizer(r)

	w := strings.Builder{}
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		t := z.Token()
		if tt == html.StartTagToken && t.Data == "body" {
			log.Info("body", "data", t.Data)

			tt = z.Next()
			for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
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
	return ioutil.NopCloser(strings.NewReader(w.String())), nil
}
