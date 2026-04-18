package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Genderize struct {
	Gender      *string `json:"gender"`
	Probability float64 `json:"probability"`
	Count       *int    `json:"count"`
}

type Agify struct {
	Age *int `json:"age"`
}

type Nationalize struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func FetchAll(name string) (*Genderize, *Agify, *Nationalize, error) {
	gURL := fmt.Sprintf("https://api.genderize.io?name=%s", name)
	aURL := fmt.Sprintf("https://api.agify.io?name=%s", name)
	nURL := fmt.Sprintf("https://api.nationalize.io?name=%s", name)

	var g Genderize
	var a Agify
	var n Nationalize

	if err := fetch(gURL, &g); err != nil {
		return nil, nil, nil, err
	}
	if err := fetch(aURL, &a); err != nil {
		return nil, nil, nil, err
	}
	if err := fetch(nURL, &n); err != nil {
		return nil, nil, nil, err
	}
	return &g, &a, &n, nil
}

func fetch(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
