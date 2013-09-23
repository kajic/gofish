package main

import (
	"flag"
	"log"

	"github.com/kajic/gofish/geonames"
)

func main() {
	addr := flag.String("a", ":8080", "listen address")
	twoFishesHost := flag.String("twoFishesHost", "localhost:8081", "two fishes server hostname (with port)")
	countryInfo := flag.String("geonames.countryInfo", "countryInfo.txt", "genonames countryinfo datafile")
	batchPoints := flag.String("batch.points", "", "path to points datafile for batch reversal")
	batchGroup := flag.Bool("batch.group", true, "applies when batch.points is given. group reversed points by country")
	flag.Parse()

	geonames.ParseCountryInfo(*countryInfo)

	if *batchPoints != "" {
		batch := NewBatch(*batchPoints, *twoFishesHost)
		batch.run(*batchGroup)
		return
	}

	log.Fatal(listenAndServe(*addr, *twoFishesHost))
}
