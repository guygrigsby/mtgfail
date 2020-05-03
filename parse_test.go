package mtgfail

import (
	"os"
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
)

func TestReadCardList(t *testing.T) {
	log := log15.New()
	f, err := os.Open("./deck.txt")
	assert.NoError(t, err)
	_, err = ReadCardList(f, log)
	assert.NoError(t, err)
}
