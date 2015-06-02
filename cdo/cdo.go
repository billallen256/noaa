package cdo

import (
	"encoding/json"
	"fmt"
	"github.com/gershwinlabs/noaa"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	BASE_URL = "http://www.ncdc.noaa.gov/cdo-web/api/v2"
)

type CDO struct {
	Metadata Metadata `json:"metadata"`
	Results  []Result `json:"results"`
}

type Metadata struct {
	Resultset Resultset `json:"resultset"`
}

type Resultset struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Result struct {
	Attributes string  `json:"attributes"`
	Datatype   string  `json:"datatype"`
	Date       string  `json:"date"`
	Station    string  `json:"station"`
	Value      float64 `json:"value"`
}

func subTimeSpans(overallTimeSpan noaa.TimeSpan) []noaa.TimeSpan {
	timeSpans := make([]noaa.TimeSpan, 0, 1)
	durationRemaining := overallTimeSpan.End.Sub(overallTimeSpan.Begin)
	begin := overallTimeSpan.Begin

	for begin.Before(overallTimeSpan.End) {
		currDuration := time.Duration(math.Min(float64(durationRemaining), float64(370*24*time.Hour)))
		end := begin.Add(currDuration)
		timeSpans = append(timeSpans, noaa.TimeSpan{begin, end})
		begin = end
		durationRemaining -= currDuration
	}

	return timeSpans
}

func FetchDataFromStationForTimeSpan(station string, overallTimeSpan noaa.TimeSpan, token string) (chan *Result, error) {
	cdoChan := make(chan *CDO)
	rateLimiter := time.Tick(1 * time.Second)
	rChan := make(chan *Result, 10)
	timeSpans := subTimeSpans(overallTimeSpan)
	logger := log.New(os.Stderr, "NOAA CDO ", log.LstdFlags)

	// goroutine 1: handle the requests and put CDO objects
	// on the channel to handle later
	go func() {
		count := 0
		offset := 1
		limit := 1000

		for _, ts := range timeSpans {
			startdate := ts.Begin.Format("2006-01-02")
			enddate := ts.End.Format("2006-01-02")

			for {
				u, _ := url.Parse(BASE_URL + "/data")

				q := u.Query()
				q.Set("datasetid", "GHCND")
				q.Set("limit", fmt.Sprintf("%d", limit))
				q.Set("stationid", station)
				q.Set("startdate", startdate)
				q.Set("enddate", enddate)
				q.Set("offset", fmt.Sprintf("%d", offset))
				q.Set("includemetadata", "true")

				u.RawQuery = q.Encode()
				req, err := http.NewRequest("GET", u.String(), nil)

				if err != nil {
					logger.Println(err)
					break
				}

				req.Header.Set("token", token)
				client := &http.Client{}
				<-rateLimiter // helps keep requests from blowing the NOAA quota
				resp, err := client.Do(req)

				if err != nil {
					logger.Println(err)
					break
				}

				if resp.StatusCode != 200 {
					logger.Println(resp.Status)
					break
				}

				var cdo CDO
				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&cdo)
				resp.Body.Close()

				if err != nil {
					logger.Println(err)
					break
				}

				count = cdo.Metadata.Resultset.Count
				logger.Printf("count=%d limit=%d offset=%d\n", count, limit, offset)
				cdoChan <- &cdo

				if count < limit+offset {
					break
				}

				offset += limit
			}

			close(cdoChan)
		}
	}()

	// goroutine 2: take individual results out of each CDO coming in
	go func() {
		for c := range cdoChan {
			for _, result := range c.Results {
				rChan <- &result
			}
		}

		close(rChan)
	}()

	return rChan, nil
}
