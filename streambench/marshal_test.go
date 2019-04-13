package streambench

import (
	"io/ioutil"
	"testing"
)

func BenchmarkMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Marshal(ioutil.Discard)
	}
}
