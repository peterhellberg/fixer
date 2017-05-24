package fixer

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestDateUnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		s    string
		want time.Time
	}{
		{`{"date":false}`, time.Time{}},
		{`{"date":"not-a-date"}`, time.Time{}},
		{`{"date":"2015-01-06"}`, time.Date(2015, 1, 6, 0, 0, 0, 0, time.UTC)},
		{`{"date":"2016-02-12"}`, time.Date(2016, 2, 12, 0, 0, 0, 0, time.UTC)},
		{`{"date":"2017-03-24"}`, time.Date(2017, 3, 24, 0, 0, 0, 0, time.UTC)},
	} {
		var v struct {
			Date Date `json:"date"`
		}

		json.NewDecoder(strings.NewReader(tt.s)).Decode(&v)

		if got := v.Date.Time; !got.Equal(tt.want) {
			t.Fatalf("v.Date.Time = %v, want %v", got, tt.want)
		}
	}
}

func TestCurrenciesString(t *testing.T) {
	for _, tt := range []struct {
		cs   Currencies
		want string
	}{
		{nil, ""},
		{Currencies{}, ""},
		{Currencies{SEK, DKK}, "DKK,SEK"},
		{Currencies{USD, AUD, EUR}, "AUD,EUR,USD"},
	} {
		if got := tt.cs.String(); got != tt.want {
			t.Fatalf("(Currencies{%s}).String() = %q, want %q", tt.cs, got, tt.want)
		}
	}
}
