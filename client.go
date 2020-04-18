package fixer

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// FixerClient is a client configured to use https://api.fixer.io
var FixerClient = NewClient(AccessKey(os.Getenv("FIXER_ACCESS_KEY")))

// ExratesClient is a client configured to use https://api.exchangeratesapi.io
var ExratesClient = NewClient(BaseURL("https://api.exchangeratesapi.io"))

// DefaultClient is the default client for the Foreign exchange rates and currency conversion API
var DefaultClient = FixerClient

// Client for the Foreign exchange rates and currency conversion API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	accessKey  string
	userAgent  string
}

// NewClient creates a Client
func NewClient(options ...func(*Client)) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseURL: &url.URL{
			Scheme: "http",
			Host:   "data.fixer.io",
			Path:   "/api",
		},
		accessKey: "",
		userAgent: "fixer/client.go (https://github.com/peterhellberg/fixer)",
	}

	for _, f := range options {
		f(c)
	}

	return c
}

// HTTPClient changes the HTTP client to the provided *http.Client
func HTTPClient(hc *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// BaseURL changes the base URL to the provided rawurl
func BaseURL(rawurl string) func(*Client) {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// AccessKey sets the access key used by the client
func AccessKey(ak string) func(*Client) {
	return func(c *Client) {
		c.accessKey = ak
	}
}

// UserAgent changes the User-Agent used by the client
func UserAgent(ua string) func(*Client) {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// Base sets the base query variable based on a Currency
func Base(c Currency) url.Values {
	v := url.Values{}

	if s := string(c); s != "" {
		v.Set("base", s)
	}

	return v
}

// Symbols sets the symbols query variable based on the provided currencies
func Symbols(cs ...Currency) url.Values {
	v := url.Values{}

	if s := Currencies(cs).String(); s != "" {
		v.Set("symbols", s)
	}

	return v
}

// Latest foreign exchange reference rates
func (c *Client) Latest(ctx context.Context, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, "/latest", c.query(attributes))
}

// At returns historical rates for any day since 1999
func (c *Client) At(ctx context.Context, t time.Time, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, "/"+c.date(t), c.query(attributes))
}

func (c *Client) date(t time.Time) string {
	return t.Format("2006-01-02")
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (*Response, error) {
	req, err := c.request(ctx, path, query)
	if err != nil {
		return nil, err
	}

	r, err := c.do(req)
	if err != nil {
		return nil, err
	}

	r.Links = Links{
		"base": c.baseURL.String(),
		"self": req.URL.String(),
	}

	return r, nil
}

func (c *Client) query(attributes []url.Values) url.Values {
	v := url.Values{}

	for _, a := range attributes {
		if base := a.Get("base"); base != "" {
			v.Set("base", base)
		}

		if symbols := a.Get("symbols"); symbols != "" {
			v.Set("symbols", symbols)
		}
	}

	return v
}

func (c *Client) request(ctx context.Context, path string, query url.Values) (*http.Request, error) {
	rawurl := c.baseURL.Path + path

	if c.accessKey != "" {
		query.Set("access_key", c.accessKey)
	}

	if len(query) > 0 {
		rawurl += "?" + query.Encode()
	}

	rel, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.baseURL.ResolveReference(rel).String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) do(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, 64)
		_ = resp.Body.Close()
	}()

	if err := responseError(resp); err != nil {
		return nil, err
	}

	var r Response

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Latest foreign exchange reference rates using the DefaultClient
func Latest(ctx context.Context, attributes ...url.Values) (*Response, error) {
	return DefaultClient.Latest(ctx, attributes...)
}

// At returns historical rates for any day since 1999 using the DefaultClient
func At(ctx context.Context, t time.Time, attributes ...url.Values) (*Response, error) {
	return DefaultClient.At(ctx, t, attributes...)
}
