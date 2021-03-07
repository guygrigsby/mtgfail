package tabletopsimulator

import (
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/require"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		normalized  string
		capitalized string
	}{
		{
			"a dog took a walk",
			"a Dog Took a Walk",
		},
		{
			"walk for the cure",
			"Walk for the Cure",
		},
	}

	log := log15.New()

	for _, tc := range tests {
		got := Capitalize(tc.normalized, log)
		require.Equal(t, tc.capitalized, got)
	}

}
