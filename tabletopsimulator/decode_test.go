package tabletopsimulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	log := log15.New()
	b, err := ioutil.ReadFile("../testdata/testdeck2.json")
	require.NoError(t, err)
	r := bytes.NewBuffer(b)
	ri, err := Decode(r, log)
	require.NoError(t, err)
	out, err := json.MarshalIndent(ri, "", "\t")
	require.NoError(t, err)
	fmt.Printf("%s", out)

}
