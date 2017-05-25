/*

Package fixer contains a client for the
Foreign exchange rates and currency conversion API

Installation

    go get -u github.com/peterhellberg/fixer

Usage

A small usage example

      package main

      import (
      	"context"
      	"flag"
      	"fmt"

      	"github.com/peterhellberg/fixer"
      )

      func main() {
      	f := flag.String("from", "EUR", "")
      	t := flag.String("to", "SEK", "")
      	n := flag.Float64("n", 1, "")

      	flag.Parse()

      	from, to := fixer.Currency(*f), fixer.Currency(*t)

      	resp, err := fixer.Latest(context.Background(),
      		fixer.Base(from), fixer.Symbols(to),
      	)

      	if err == nil {
      		fmt.Printf("%.2f %s equals %.2f %s\n", *n, from, resp.Rates[to]**n, to)
      	}
      }

API Documentation

http://fixer.io/

*/
package fixer

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

// Date wraps time.Time
type Date struct {
	time.Time
}

// UnmarshalJSON parses dates in YYYY-MM-DD format
func (d *Date) UnmarshalJSON(b []byte) error {
	var value string

	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	if err != nil {
		return err
	}

	*d = Date{t}

	return nil
}

// Rates is the list of rates quoted against the base (EUR by default)
type Rates map[Currency]float64

// Links is a links object related to the primary data of the Response
type Links map[string]string

// Response data from the Foreign exchange rates and currency conversion API
type Response struct {
	Base  Currency `json:"base"`
	Date  Date     `json:"date"`
	Rates Rates    `json:"rates"`
	Links Links    `json:"links,omitempty"`
}

// Currencies is a slice of Currency
type Currencies []Currency

func (cs Currencies) String() string {
	symbols := []string{}

	for _, c := range cs {
		symbols = append(symbols, string(c))
	}

	sort.Strings(symbols)

	return strings.Join(symbols, ",")
}

// Currency is the type used for ISO 4217 Currency codes
type Currency string

// Currency codes published by the European Central Bank
const (
	AUD Currency = "AUD"
	BGN Currency = "BGN"
	BRL Currency = "BRL"
	CAD Currency = "CAD"
	CHF Currency = "CHF"
	CNY Currency = "CNY"
	CZK Currency = "CZK"
	DKK Currency = "DKK"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	HKD Currency = "HKD"
	HRK Currency = "HRK"
	HUF Currency = "HUF"
	IDR Currency = "IDR"
	ILS Currency = "ILS"
	INR Currency = "INR"
	JPY Currency = "JPY"
	KRW Currency = "KRW"
	MXN Currency = "MXN"
	MYR Currency = "MYR"
	NOK Currency = "NOK"
	NZD Currency = "NZD"
	PHP Currency = "PHP"
	PLN Currency = "PLN"
	RON Currency = "RON"
	RUB Currency = "RUB"
	SEK Currency = "SEK"
	SGD Currency = "SGD"
	THB Currency = "THB"
	TRY Currency = "TRY"
	USD Currency = "USD"
	ZAR Currency = "ZAR"
)
