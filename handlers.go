package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

var (
	templates = template.Must(template.New("").
		Funcs(template.FuncMap{
			"renderMoney": renderMoney,
		}).ParseGlob("templates/*.html"))
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	l := hlog.FromRequest(r)
	curCurr := currentCurrency(r)

	currencies := Currencies()
	products := ListProducts()

	l.Info().Str("currency", curCurr).Int("cur num", len(currencies)).Int("prod num", len(products)).Msg("home handler")

	type productView struct {
		Item  Product
		Price Money
	}
	ps := make([]productView, len(products))
	for i, p := range products {
		price := Convert(p.PriceUsd, curCurr)
		ps[i] = productView{p, price}
	}

	rid, _ := hlog.IDFromRequest(r)
	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
		"request_id":    rid.String(),
		"user_currency": curCurr,
		"currencies":    currencies,
		"products":      ps,
		//"cart_size":     len(cart),
		//"banner_color":  os.Getenv("BANNER_COLOR"), // illustrates canary deployments
	}); err != nil {
		log.Info().Err(err).Msg("unable to parse home template")
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	l := hlog.FromRequest(r)
	id := chi.URLParam(r, "id")
	if id == "" {
		renderHTTPError(l, r, w, errors.New("product id not specified"), http.StatusBadRequest)
		return
	}
	l.Debug().Str("id", id).Str("currency", currentCurrency(r)).Msg("serving product page") //

	p, err := GetProduct(id)
	if err != nil {
		renderHTTPError(l, r, w, errors.Wrap(err, "could not retrieve product"), http.StatusInternalServerError)
		return
	}

	currencies := Currencies()
	price := Convert(p.PriceUsd, currentCurrency(r))
	product := struct {
		Item  Product
		Price Money
	}{*p, price}

	rid, _ := hlog.IDFromRequest(r)
	if err := templates.ExecuteTemplate(w, "product", map[string]interface{}{
		"request_id":    rid.String(),
		"user_currency": currentCurrency(r),
		"currencies":    currencies,
		"product":       product,
	}); err != nil {
		l.Info().Err(err).Msg("unable to parse product template")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	l := hlog.FromRequest(r)
	l.Debug().Msg("logging out")
	for _, c := range r.Cookies() {
		c.Expires = time.Now().Add(-time.Hour * 24 * 365)
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

func setCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	l := hlog.FromRequest(r)
	cur := r.FormValue("currency_code")
	l.Info().Str("curr.new", cur).Str("curr.old", currentCurrency(r)).Msg("setting currency") //

	if cur != "" {
		http.SetCookie(w, &http.Cookie{
			Name:   cookieCurrency,
			Value:  cur,
			MaxAge: cookieMaxAge,
		})
	}
	referer := r.Header.Get("referer")
	if referer == "" {
		referer = "/"
	}
	w.Header().Set("Location", referer)
	w.WriteHeader(http.StatusFound)
}

func renderHTTPError(l *zerolog.Logger, r *http.Request, w http.ResponseWriter, err error, code int) {
	l.Error().Err(err).Msg("request error")
	errMsg := fmt.Sprintf("%+v", err)

	w.WriteHeader(code)
	rid, _ := hlog.IDFromRequest(r)
	templates.ExecuteTemplate(w, "error", map[string]interface{}{
		"request_id":  rid.String(),
		"error":       errMsg,
		"status_code": code,
		"status":      http.StatusText(code)})
}

func currentCurrency(r *http.Request) string {
	c, _ := r.Cookie(cookieCurrency)
	if c != nil {
		return c.Value
	}
	return defaultCurrency
}

func renderMoney(money Money) string {
	return fmt.Sprintf("%s %d.%02d", money.CurrencyCode, money.Units, money.Nanos/10000000)
}
