package handler

import (
	"unicode"
)

// ValueFetcher accept two kind of values:
// 1. lower camel case ex: externalKey
// 2. pascal case ex: external_key
// it will auto detect input v's type and trying to fetch value from f function in those naming rule.
func ValueFetcher(v string, f func(key string) string) string {
	if f(v) != "" {
		return f(v)
	}

	nv := ""
	isConvert := false
	for _, s := range v {
		if unicode.IsUpper(s) {
			nv += "_" + string(unicode.ToLower(s))
		} else if s == '_' {
			isConvert = true
		} else {
			if isConvert {
				nv += string(unicode.ToUpper(s))
				isConvert = false
			} else {
				nv += string(s)
			}
		}
	}

	return f(nv)
}
