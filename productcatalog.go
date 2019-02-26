package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

var (
	prodList []Product
)

func init() {
	c, err := ioutil.ReadFile("products.json")
	if err != nil {
		log.Fatalf("failed to open product catalog json file: %v", err)
	}
	pl := map[string][]Product{}
	if err := json.Unmarshal(c, &pl); err != nil {
		log.Fatalf("failed to parse the catalog JSON: %v", err)
	}
	prodList = pl["products"]
	log.Printf("successfully parsed %d products catalog from json\n", len(prodList))
}

// Product
type Product struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Picture     string `json:"picture,omitempty"`
	PriceUsd    Money  `json:"priceUsd,omitempty"`
	// Categories such as "vintage" or "gardening" that can be used to look up
	// other related products.
	Categories []string `json:"categories,omitempty"`
}

func ListProducts() []Product {
	return prodList
}

func GetProduct(pid string) (*Product, error) {
	var found *Product
	for i := 0; i < len(prodList); i++ {
		if pid == prodList[i].Id {
			found = &prodList[i]
		}
	}
	if found == nil {
		return nil, errors.New("no product with ID " + pid)
	}
	return found, nil
}

func SearchProducts(query string) ([]Product, error) {
	// Intepret query as a substring match in name or description.
	var ps []Product
	for _, p := range prodList {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(query)) {
			ps = append(ps, p)
		}
	}
	return ps, nil
}
