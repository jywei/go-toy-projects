package streambench

import (
	"io/ioutil"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Encode(ioutil.Discard)
	}
}
