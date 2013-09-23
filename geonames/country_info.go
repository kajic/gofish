package geonames

import (
	"bufio"
	"os"
	"strings"
)

var countries = map[string]*Country{}

type Country struct {
	CountryCode string `json:"cc"`
	Name        string `json:"name"`
}

func NewCountry(line string) *Country {
	cols := strings.Split(line, "\t")

	return &Country{
		CountryCode: cols[0],
		Name:        cols[4],
	}
}

func ParseCountryInfo(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == "#"[0] {
			continue
		}

		country := NewCountry(line)
		countries[country.CountryCode] = country
	}

	return nil
}

func GetCountry(cc string) *Country {
	return countries[cc]
}
