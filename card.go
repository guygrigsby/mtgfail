package mtgfail

import "encoding/base64"

func Key(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}
