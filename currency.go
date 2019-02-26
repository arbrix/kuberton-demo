package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
)

var rates = map[string]float64{}

type xmlCurRate struct {
	XMLName xml.Name `xml:"Cube"`
	Cur     string   `xml:"currency,attr"`
	Rate    string   `xml:"rate,attr"`
}

type xmlCube struct {
	XMLName xml.Name     `xml:"Cube"`
	Rates   []xmlCurRate `xml:"Cube"`
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
	for _, cr := range x.Cube.Rates {
		r, err = strconv.ParseFloat(cr.Rate, 64)
		if err != nil {
			continue
		}
		rates[cr.Cur] = r
	}

	rates["EUR"] = 1.0
}

func Rates() map[string]float64 {
	return rates
}
