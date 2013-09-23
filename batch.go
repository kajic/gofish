package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kajic/gofish/geonames"
)

type CountryGroup struct {
	Country *geonames.Country
	Count   int
}

type CountryStatsMap map[string]*CountryGroup

func (stats CountryStatsMap) Incr(country *geonames.Country) {
	countryStats, ok := stats[country.CountryCode]
	if !ok {
		stats[country.CountryCode] = &CountryGroup{country, 1}
	} else {
		countryStats.Count++
	}
}

type CountryStatsList []*CountryGroup

func (stats CountryStatsList) Len() int {
	return len(stats)
}
func (stats CountryStatsList) Less(i, j int) bool {
	return stats[i].Count < stats[j].Count
}
func (stats CountryStatsList) Swap(i, j int) {
	stats[i], stats[j] = stats[j], stats[i]
}

type ReverseResult struct {
	country *geonames.Country
	lat     float64
	lng     float64
}

type Batch struct {
	path          string
	twoFishesHost string
	countryc      chan *ReverseResult
	errc          chan error
	semc          chan bool
}

func NewBatch(path, twoFishesHost string) *Batch {
	return &Batch{
		path,
		twoFishesHost,
		make(chan *ReverseResult),
		make(chan error),
		make(chan bool, 16),
	}
}

func (me *Batch) reverse(lat float64, lng float64) {
	country, err := reverse(me.twoFishesHost, lat, lng)
	reverseResult := &ReverseResult{country, lat, lng}
	<-me.semc
	if err != nil {
		me.errc <- err
	} else {
		me.countryc <- reverseResult
	}
}

func (me *Batch) enqueue() (int, error) {
	file, err := os.Open(me.path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	jobsCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		llParts := strings.Split(line, "\t")
		lat, err := strconv.ParseFloat(llParts[0], 64)
		if err != nil {
			me.errc <- err
			continue
		}
		lng, err := strconv.ParseFloat(llParts[1], 64)
		if err != nil {
			me.errc <- err
			continue
		}

		me.semc <- true
		go me.reverse(lat, lng)
		jobsCount++
	}
	return jobsCount, nil
}

func (me *Batch) collect(jobsCount int, callback func(*ReverseResult)) {
Perform:
	for ; jobsCount > 0; jobsCount-- {
		select {
		case reverseResult, ok := <-me.countryc:
			if !ok {
				break Perform
			}
			callback(reverseResult)
		case err := <-me.errc:
			log.Println(err)
		}
	}
}

func (me *Batch) stats() (CountryStatsList, error) {
	jobsCount, err := me.enqueue()
	if err != nil {
		return nil, err
	}
	countryStats := CountryStatsMap{}
	me.collect(jobsCount, func(reverseResult *ReverseResult) {
		countryStats.Incr(reverseResult.country)
	})
	countryStatsList := CountryStatsList{}
	for _, value := range countryStats {
		countryStatsList = append(countryStatsList, value)
	}
	sort.Sort(sort.Reverse(countryStatsList))
	return countryStatsList, nil
}

func (me *Batch) plain() error {
	jobsCount, err := me.enqueue()
	if err != nil {
		return err
	}
	me.collect(jobsCount, func(reverseResult *ReverseResult) {
		fmt.Println(reverseResult.lat, reverseResult.lng, reverseResult.country.Name)
	})
	return nil
}

func (me *Batch) run(group bool) {
	t := time.Now()

	if group {
		stats, err := me.stats()
		if err != nil {
			fmt.Println(err)
		}
		for _, countryCount := range stats {
			fmt.Println(countryCount.Count, countryCount.Country.Name)
		}
	} else {
		me.plain()
	}
	fmt.Println("time", time.Now().Sub(t))
}
