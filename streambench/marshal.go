package streambench

import (
	"encoding/json"
	"io"
)

// Marshal does some marshaling for testing
func Marshal(w io.Writer) {
	j, _ := json.Marshal(&Data{
		ID:       "abcde",
		Name:     "Roy",
		Email:    "roywjy@gmail.com",
		Whatever: 100,
	})
	w.Write(j)
}
