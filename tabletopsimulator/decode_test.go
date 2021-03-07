package tabletopsimulator

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	b, err := ioutil.ReadFile("../testdata/AlelaTTSFormat.json")
	require.NoError(t, err)
	r := bytes.NewBuffer(b)
	ri, err := Decode(r)
	require.NoError(t, err)
	err

}
