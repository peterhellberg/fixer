# fixer

[![Build Status](https://travis-ci.org/peterhellberg/fixer.svg?branch=master)](https://travis-ci.org/peterhellberg/fixer)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/fixer)](https://goreportcard.com/report/github.com/peterhellberg/fixer)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/fixer)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/fixer#license-mit)

Go client for [Fixer.io](http://fixer.io/) (Foreign exchange rates and currency conversion API)

## Installation

    go get -u github.com/peterhellberg/fixer

## Usage examples

**SEK quoted against USD and EUR**

```go
package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/peterhellberg/fixer"
)

func main() {
	f := fixer.NewClient()

	resp, err := f.Latest(context.Background(),
		fixer.Base(fixer.SEK),
		fixer.Symbols(
			fixer.USD,
			fixer.EUR,
		),
	)
	if err != nil {
		return
	}

	encode(resp)
}

func encode(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	enc.Encode(v)
}
```

```json
{
 "base": "SEK",
 "date": "2017-05-24T00:00:00Z",
 "rates": {
  "EUR": 0.10265,
  "USD": 0.1149
 },
 "links": {
  "base": "https://api.fixer.io",
  "self": "https://api.fixer.io/latest?base=SEK&symbols=EUR%2CUSD"
 }
}
```

**Using [goexrates](http://goexrates.mikolajczakluq.com/) instead**

```go
package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/peterhellberg/fixer"
)

func main() {
	f := fixer.NewClient(
		fixer.BaseURL("http://exr.mikolajczakluq.com"),
	)

	resp, err := f.At(context.Background(), time.Now(),
		fixer.Base(fixer.GBP),
		fixer.Symbols(
			fixer.SEK,
			fixer.NOK,
		),
	)
	if err != nil {
		return
	}

	encode(resp)
}

func encode(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	enc.Encode(v)
}
```

```json
{
 "base": "GBP",
 "date": "2017-05-24T00:00:00Z",
 "rates": {
  "NOK": 10.86901,
  "SEK": 11.28307
 },
 "links": {
  "base": "http://exr.mikolajczakluq.com",
  "self": "http://exr.mikolajczakluq.com/2017-05-25?base=GBP&symbols=NOK%2CSEK"
 }
}
```

## API documentation

<http://fixer.io/>

## License (MIT)

Copyright (c) 2017 [Peter Hellberg](https://c7.se/)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the "Software"),
> to deal in the Software without restriction, including without limitation
> the rights to use, copy, modify, merge, publish, distribute, sublicense,
> and/or sell copies of the Software, and to permit persons to whom the
> Software is furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included
> in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
> OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
> IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
> DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
> TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
> OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
