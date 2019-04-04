package handler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/honestbee/habitat-catalog-cache/handler"
)

func TestValueFetcher(t *testing.T) {
	testCases := []struct {
		value string
		fn    func(string) string
		want  string
	}{
		{
			value: "externalKey",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return key
				}
				return ""
			},
			want: "external_key",
		},
		{
			value: "external_key",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return ""
				}
				return key
			},
			want: "externalKey",
		},
		{
			value: "external_key",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return key
				}
				return ""
			},
			want: "external_key",
		},
		{
			value: "",
			fn: func(key string) string {
				return ""
			},
			want: "",
		},
		{
			value: "externalKeY",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return key
				}
				return ""
			},
			want: "external_ke_y",
		},
		{
			value: "external_key_",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return ""
				}
				return key
			},
			want: "externalKey",
		},
		{
			value: "this_is_a_love_song",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return ""
				}
				return key
			},
			want: "thisIsALoveSong",
		},
		{
			value: "thisIsALoveSong",
			fn: func(key string) string {
				if strings.Contains(key, "_") {
					return key
				}
				return ""
			},
			want: "this_is_a_love_song",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("value:%s, want:%s", tc.value, tc.want), func(t *testing.T) {
			got := handler.ValueFetcher(tc.value, tc.fn)
			if tc.want != got {
				t.Errorf("want:%s, got:%s", tc.want, got)
			}
		})
	}
}
