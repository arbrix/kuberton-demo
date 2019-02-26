package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var rates = map[string]float64{}

type xmlCurRate struct {
	XMLName xml.Name `xml:"Cube"`
	Cur     string   `xml:"currency,attr"`
	Rate    string   `xml:"rate,attr"`
}

type xmlCube1 struct {
	XMLName xml.Name     `xml:"Cube"`
	Rates   []xmlCurRate `xml:"Cube"`
}

type xmlCube struct {
	XMLName xml.Name `xml:"Cube"`
	Cube    xmlCube1 `xml:"Cube"`
}

type xmlEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Cube    xmlCube  `xml:"Cube"`
}

const (
	urlSrc          = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	defaultCurrency = "USD"

	cookieMaxAge = 60 * 60 * 48

	cookiePrefix   = "shop_"
	cookieCurrency = cookiePrefix + "currency"
)

var whitelistedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"CAD": true,
	"JPY": true,
	"GBP": true,
	"TRY": true}

func init() {
	res, err := http.Get(urlSrc)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatalf("unable to request rates: code: %d, err: %v", res.StatusCode, err)
	}

	var x xmlEnvelope
	if err = xml.NewDecoder(res.Body).Decode(&x); err != nil {
		log.Fatalf("unable to parse currency responce: %v", err)
	}

	var r float64
	for _, cr := range x.Cube.Cube.Rates {
		r, err = strconv.ParseFloat(cr.Rate, 64)
		if err != nil || !whitelistedCurrencies[cr.Cur] {
			continue
		}
		rates[cr.Cur] = r
	}

	rates["EUR"] = 1.0
	log.Printf("currencies rates successfully retraived: %d\n", len(rates))
}

func Rates() map[string]float64 {
	return rates
}

func Currencies() []string {
	cs := []string{}
	for c := range rates {
		if whitelistedCurrencies[c] {
			cs = append(cs, c)
		}
	}
	sort.Strings(cs)
	return cs
}

func Convert(price Money, currency string) Money {
	if currency == "USD" {
		return price
	}

	if rates[price.CurrencyCode] == 0.0 || rates[currency] == 0.0 {
		return Money{CurrencyCode: currency}
	}

	eurUnits := float64(price.Units) / rates[price.CurrencyCode]
	eurNanos := float64(price.Nanos) / rates[price.CurrencyCode]
	return Money{
		CurrencyCode: currency,
		Units:        int64(eurUnits * rates[currency]),
		Nanos:        int32(eurNanos * rates[currency]),
	}
}
