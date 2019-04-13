package streambench

import (
	"encoding/json"
	"io"
)

// Encode encodes something for testing
func Encode(w io.Writer) {
	json.NewEncoder(w).Encode(&Data{
		ID:       "abcde",
		Name:     "Roy",
		Email:    "roywjy@gmail.com",
		Whatever: 100,
	})
}
