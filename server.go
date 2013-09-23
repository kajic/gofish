package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kajic/gofish/geonames"
)

type Response struct {
	Country *geonames.Country `json:"country"`
	Error   string            `json:"error"`
}

func getLLFromUrl(u *url.URL) (float64, float64, error) {
	ll := u.Query().Get("ll")
	if ll == "" {
		return 0, 0, errors.New("The ll parameter is empty")
	}
	llParts := strings.Split(ll, ",")
	if len(llParts) != 2 {
		return 0, 0, errors.New("The ll parameter requires two comma separated floats")
	}
	lat, err := strconv.ParseFloat(llParts[0], 64)
	if err != nil {
		return 0, 0, err
	}
	lng, err := strconv.ParseFloat(llParts[1], 64)
	if err != nil {
		return 0, 0, err
	}
	return lat, lng, nil
}

func jsonResponse(w io.Writer, resp interface{}) {
	jsonb, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, string(jsonb))
	}
}

type Reverse struct {
	twoFishesHost string
}

func (me *Reverse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lat, lng, err := getLLFromUrl(r.URL)
	if err != nil {
		jsonResponse(w, Response{Error: err.Error()})
		return
	}

	country, err := reverse(me.twoFishesHost, lat, lng)
	if err != nil {
		jsonResponse(w, Response{Error: err.Error()})
		return
	}

	jsonResponse(w, Response{Country: country})
}

func listenAndServe(addr, twoFishesHost string) error {
	http.Handle("/reverse", &Reverse{twoFishesHost})
	log.Printf("server listening on %s..", addr)
	return http.ListenAndServe(addr, nil)
}
