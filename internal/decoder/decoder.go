package decoder

import "github.com/gorilla/schema"

var Decoder = schema.NewDecoder()

func init() {
	Decoder.IgnoreUnknownKeys(true)
	Decoder.SetAliasTag("query")
}
