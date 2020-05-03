package mtgfail

import (
	"fmt"
	"testing"

	"github.com/inconshreveable/log15"
)

func TestReadCardList(t *testing.T) {
	log := log15.New()
	files, err := ReadCardList("./list.txt", log)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", files)
}
