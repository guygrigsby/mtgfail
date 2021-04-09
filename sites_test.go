package mtgfail

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranslateDomain(t *testing.T) {
	tc := []struct {
		name        string
		uri         string
		returnURI   *url.URL
		expectError bool
	}{
		{
			"betz",
			"https://tappedout.net/mtg-decks/betz-rakdos/",
			must("https://tappedout.net/mtg-decks/betz-rakdos/?fmt=txt"),
			false,
		},
		{
			"deckbox valid",

			"https://deckbox.org/sets/2805419",
			must("https://deckbox.org/sets/2805419/export"),
			false,
		},
		{
			"unsupported site domain (for now ;))",
			"https://scryfall.org/sets/2805419",
			must(""),
			true,
		},
		{
			"no domain",
			"/sets/2805419",
			must(""),
			true,
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			returnURI, err := TranslateDomain(test.uri)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equalf(t, test.returnURI, returnURI, "incorrect fullURL value: expected '%s', got '%s'", test.returnURI, returnURI)
			}
			t.Log("Error message", err)
		})
	}
}

func must(u string) *url.URL {
	uri, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return uri
}
