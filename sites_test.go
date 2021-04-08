package mtgfail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDomain(t *testing.T) {
	tc := []struct {
		uri      string
		valid    bool
		fullURI  string
		hasError bool
	}{
		{
			"https://tappedout.net/mtg-decks/betz-rakdos/",
			true,
			"https://tappedout.net/mtg-decks/betz-rakdos/?fmt=txt",
			false,
		},
		{
			"https://deckbox.org/sets/2805419",
			true,
			"https://deckbox.org/sets/2805419/export",
			false,
		},
		{
			"https://scryfall.org/sets/2805419",
			false,
			"",
			false,
		},
		{
			"/sets/2805419",
			false,
			"",
			true,
		},
		//{
		//	"http://805419",
		//	false,
		//	"",
		//	true,
		//},
	}

	for _, test := range tc {
		test := test
		t.Run(test.uri, func(t *testing.T) {
			t.Parallel()
			valid, fullURI, err := ParseDomain(test.uri)
			require.Equalf(t, test.hasError, err != nil,
				"incorrect error return : expected %v, got %v", test.hasError, err != nil)
			require.Equalf(t, test.fullURI, fullURI, "incorrect fullURL value: expected %s, got %s", test.fullURI, fullURI)
			require.Equalf(t, test.valid, valid, "incorrect valid value: expected %v, got %v valid")
		})
	}
}
