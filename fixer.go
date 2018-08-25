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
	AED Currency = "AED"
	AFN Currency = "AFN"
	ALL Currency = "ALL"
	AMD Currency = "AMD"
	ANG Currency = "ANG"
	AOA Currency = "AOA"
	ARS Currency = "ARS"
	AUD Currency = "AUD"
	AWG Currency = "AWG"
	AZN Currency = "AZN"
	BAM Currency = "BAM"
	BBD Currency = "BBD"
	BDT Currency = "BDT"
	BGN Currency = "BGN"
	BHD Currency = "BHD"
	BIF Currency = "BIF"
	BMD Currency = "BMD"
	BND Currency = "BND"
	BOB Currency = "BOB"
	BRL Currency = "BRL"
	BSD Currency = "BSD"
	BTC Currency = "BTC"
	BTN Currency = "BTN"
	BWP Currency = "BWP"
	BYN Currency = "BYN"
	BYR Currency = "BYR"
	BZD Currency = "BZD"
	CAD Currency = "CAD"
	CDF Currency = "CDF"
	CHF Currency = "CHF"
	CLF Currency = "CLF"
	CLP Currency = "CLP"
	CNY Currency = "CNY"
	COP Currency = "COP"
	CRC Currency = "CRC"
	CUC Currency = "CUC"
	CUP Currency = "CUP"
	CVE Currency = "CVE"
	CZK Currency = "CZK"
	DJF Currency = "DJF"
	DKK Currency = "DKK"
	DOP Currency = "DOP"
	DZD Currency = "DZD"
	EGP Currency = "EGP"
	ERN Currency = "ERN"
	ETB Currency = "ETB"
	EUR Currency = "EUR"
	FJD Currency = "FJD"
	FKP Currency = "FKP"
	GBP Currency = "GBP"
	GEL Currency = "GEL"
	GGP Currency = "GGP"
	GHS Currency = "GHS"
	GIP Currency = "GIP"
	GMD Currency = "GMD"
	GNF Currency = "GNF"
	GTQ Currency = "GTQ"
	GYD Currency = "GYD"
	HKD Currency = "HKD"
	HNL Currency = "HNL"
	HRK Currency = "HRK"
	HTG Currency = "HTG"
	HUF Currency = "HUF"
	IDR Currency = "IDR"
	ILS Currency = "ILS"
	IMP Currency = "IMP"
	INR Currency = "INR"
	IQD Currency = "IQD"
	IRR Currency = "IRR"
	ISK Currency = "ISK"
	JEP Currency = "JEP"
	JMD Currency = "JMD"
	JOD Currency = "JOD"
	JPY Currency = "JPY"
	KES Currency = "KES"
	KGS Currency = "KGS"
	KHR Currency = "KHR"
	KMF Currency = "KMF"
	KPW Currency = "KPW"
	KRW Currency = "KRW"
	KWD Currency = "KWD"
	KYD Currency = "KYD"
	KZT Currency = "KZT"
	LAK Currency = "LAK"
	LBP Currency = "LBP"
	LKR Currency = "LKR"
	LRD Currency = "LRD"
	LSL Currency = "LSL"
	LTL Currency = "LTL"
	LVL Currency = "LVL"
	LYD Currency = "LYD"
	MAD Currency = "MAD"
	MDL Currency = "MDL"
	MGA Currency = "MGA"
	MKD Currency = "MKD"
	MMK Currency = "MMK"
	MNT Currency = "MNT"
	MOP Currency = "MOP"
	MRO Currency = "MRO"
	MUR Currency = "MUR"
	MVR Currency = "MVR"
	MWK Currency = "MWK"
	MXN Currency = "MXN"
	MYR Currency = "MYR"
	MZN Currency = "MZN"
	NAD Currency = "NAD"
	NGN Currency = "NGN"
	NIO Currency = "NIO"
	NOK Currency = "NOK"
	NPR Currency = "NPR"
	NZD Currency = "NZD"
	OMR Currency = "OMR"
	PAB Currency = "PAB"
	PEN Currency = "PEN"
	PGK Currency = "PGK"
	PHP Currency = "PHP"
	PKR Currency = "PKR"
	PLN Currency = "PLN"
	PYG Currency = "PYG"
	QAR Currency = "QAR"
	RON Currency = "RON"
	RSD Currency = "RSD"
	RUB Currency = "RUB"
	RWF Currency = "RWF"
	SAR Currency = "SAR"
	SBD Currency = "SBD"
	SCR Currency = "SCR"
	SDG Currency = "SDG"
	SEK Currency = "SEK"
	SGD Currency = "SGD"
	SHP Currency = "SHP"
	SLL Currency = "SLL"
	SOS Currency = "SOS"
	SRD Currency = "SRD"
	STD Currency = "STD"
	SVC Currency = "SVC"
	SYP Currency = "SYP"
	SZL Currency = "SZL"
	THB Currency = "THB"
	TJS Currency = "TJS"
	TMT Currency = "TMT"
	TND Currency = "TND"
	TOP Currency = "TOP"
	TRY Currency = "TRY"
	TTD Currency = "TTD"
	TWD Currency = "TWD"
	TZS Currency = "TZS"
	UAH Currency = "UAH"
	UGX Currency = "UGX"
	USD Currency = "USD"
	UYU Currency = "UYU"
	UZS Currency = "UZS"
	VEF Currency = "VEF"
	VND Currency = "VND"
	VUV Currency = "VUV"
	WST Currency = "WST"
	XAF Currency = "XAF"
	XAG Currency = "XAG"
	XAU Currency = "XAU"
	XCD Currency = "XCD"
	XDR Currency = "XDR"
	XOF Currency = "XOF"
	XPF Currency = "XPF"
	YER Currency = "YER"
	ZAR Currency = "ZAR"
	ZMK Currency = "ZMK"
	ZMW Currency = "ZMW"
	ZWL Currency = "ZWL"
)
