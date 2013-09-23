package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kajic/gofish/geonames"
)

var twoFishesUrl = "http://%s/?responseIncludes=PARENTS&&maxInterpretations=1&ll=%f,%f"

type TwoFishesFeature struct {
	CountryCode string `json:"cc"`
}
type TwoFishesInterpretation struct {
	Feature TwoFishesFeature `json:"feature"`
}
type TwoFishesResult struct {
	Interpretations []TwoFishesInterpretation `json:"interpretations"`
}

func (tfr *TwoFishesResult) Country() (*geonames.Country, error) {
	if len(tfr.Interpretations) < 1 {
		return nil, errors.New("There are no results")
	}
	cc := tfr.Interpretations[0].Feature.CountryCode
	return geonames.GetCountry(cc), nil
}

func reverse(twoFishesHost string, lat float64, lng float64) (*geonames.Country, error) {
	u := fmt.Sprintf(twoFishesUrl, twoFishesHost, lat, lng)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	twoFishesResult := &TwoFishesResult{}
	err = json.Unmarshal(body, twoFishesResult)
	if err != nil {
		return nil, err
	}

	return twoFishesResult.Country()
}
