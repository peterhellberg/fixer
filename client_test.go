package fixer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := NewClient()

		if got, want := c.httpClient.Timeout, 20*time.Second; got != want {
			t.Fatalf("c.httpClient.Timeout = %q, want %q", got, want)
		}

		if got, want := c.baseURL.String(), "https://api.fixer.io"; got != want {
			t.Fatalf("c.baseURL.String() = %q, want %q", got, want)
		}

		if got, want := c.userAgent, "fixer/client.go (https://github.com/peterhellberg/fixer)"; got != want {
			t.Fatalf("c.userAgent = %q, want %q", got, want)
		}
	})

	t.Run("HTTPClient", func(t *testing.T) {
		c := NewClient(HTTPClient(&http.Client{
			Timeout: 40 * time.Second,
		}))

		if got, want := c.httpClient.Timeout, 40*time.Second; got != want {
			t.Fatalf("c.httpClient.Timeout = %q, want %q", got, want)
		}
	})

	t.Run("BaseURL", func(t *testing.T) {
		rawurl := "http://exr.mikolajczakluq.com"

		c := NewClient(BaseURL(rawurl))

		if got, want := c.baseURL.String(), rawurl; got != want {
			t.Fatalf("c.baseURL.String() = %q, want %q", got, want)
		}
	})

	t.Run("UserAgent", func(t *testing.T) {
		ua := "Custom User-Agent"

		c := NewClient(UserAgent(ua))

		if got, want := c.userAgent, ua; got != want {
			t.Fatalf("c.userAgent = %q, want %q", got, want)
		}
	})
}

func TestBase(t *testing.T) {
	for _, tt := range []struct {
		c    Currency
		want string
	}{
		{"", ""},
		{SEK, "base=SEK"},
		{EUR, "base=EUR"},
	} {
		if got := Base(tt.c).Encode(); got != tt.want {
			t.Fatalf("Base(%v).Encode() = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestSymbols(t *testing.T) {
	for _, tt := range []struct {
		currencies Currencies
		want       string
	}{
		{nil, ""},
		{Currencies{}, ""},
		{Currencies{SEK}, "symbols=SEK"},
		{Currencies{SEK, DKK}, "symbols=DKK%2CSEK"},
		{Currencies{EUR, USD, RUB}, "symbols=EUR%2CRUB%2CUSD"},
		{Currencies{CAD, BGN, AUD}, "symbols=AUD%2CBGN%2CCAD"},
	} {
		if got := Symbols(tt.currencies...).Encode(); got != tt.want {
			t.Fatalf("Symbols(%s).Encode() = %q, want %q", tt.currencies, got, tt.want)
		}
	}
}

func TestLatest(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	y, m, d := time.Now().Date()

	t.Run("default", func(t *testing.T) {
		resp, err := c.Latest(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := resp.Base, EUR; got != want {
			t.Fatalf("resp.Base = %q, want %q", got, want)
		}

		if got, want := resp.Date.Time, time.Date(y, m, d, 0, 0, 0, 0, time.UTC); got != want {
			t.Fatalf("resp.Date.Time = %v, want %v", got, want)
		}
	})

	t.Run("base-SEK-symbols-USD-GBP", func(t *testing.T) {
		resp, err := c.Latest(context.Background(), Base(SEK), Symbols(USD, GBP))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := resp.Base, SEK; got != want {
			t.Fatalf("resp.Base = %q, want %q", got, want)
		}

		if got, want := resp.Date.Time, time.Date(y, m, d, 0, 0, 0, 0, time.UTC); got != want {
			t.Fatalf("resp.Date.Time = %v, want %v", got, want)
		}

		if got, want := len(resp.Rates), 2; got != want {
			t.Fatalf("len(%v) = %d, want %d", resp.Rates, got, want)
		}

		if got, want := resp.Rates[GBP], 0.088628; got != want {
			t.Fatalf("resp.Rates[GBP] = %f, want %f", got, want)
		}
	})
}

func TestAt(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	t.Run("2012-03-28", func(t *testing.T) {
		date := time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)

		resp, err := c.At(context.Background(), date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := resp.Base, EUR; got != want {
			t.Fatalf("resp.Base = %q, want %q", got, want)
		}

		if got, want := resp.Date.Time, date; got != want {
			t.Fatalf("resp.Date.Time = %v, want %v", got, want)
		}
	})

	t.Run("too-old", func(t *testing.T) {
		date := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)

		if _, err := c.At(context.Background(), date); err != ErrUnprocessableEntity {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not-json", func(t *testing.T) {
		date := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

		_, err := c.At(context.Background(), date)

		if err == nil {
			t.Fatalf("expected to get error")
		}

		if got, want := err.Error(), "invalid character 'N' looking for beginning of value"; got != want {
			t.Fatalf("err.Error() = %q, want %q", got, want)
		}
	})
}

func TestGet(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	t.Run("invalid-request", func(t *testing.T) {
		_, err := c.get(context.Background(), ":/", url.Values{})

		if err == nil {
			t.Fatalf("expected to get error")
		}

		if got, want := err.Error(), "parse :/: missing protocol scheme"; got != want {
			t.Fatalf("err.Error() = %q, want %q", got, want)
		}
	})
}

func testServerAndClient() (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			enc := json.NewEncoder(w)

			w.Header().Set("Content-Type", "application/json")

			switch r.URL.String() {
			case "/latest":
				enc.Encode(map[string]interface{}{
					"base": EUR,
					"date": time.Now().Format("2006-01-02"),
				})
			case "/latest?base=SEK&symbols=GBP%2CUSD":
				enc.Encode(map[string]interface{}{
					"base": SEK,
					"date": time.Now().Format("2006-01-02"),
					"rates": Rates{
						GBP: 0.088628,
						USD: 0.1149,
					},
				})
			case "/2012-03-28":
				enc.Encode(map[string]interface{}{
					"base": EUR,
					"date": "2012-03-28",
				})
			case "/1999-12-31":
				w.WriteHeader(http.StatusUnprocessableEntity)
			case "/2000-01-01":
				w.Write([]byte(`NOT JSON`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))

	return ts, NewClient(BaseURL(ts.URL))
}
