package cdo

import (
	"encoding/json"
	"github.com/gershwinlabs/noaa"
	"net/http"
	"net/url"
)

const (
	BASE_URL = "http://www.ncdc.noaa.gov/cdo-web/api/v2"
)

type CDO struct {
	Metadata Metadata `json:"metadata"`
	Results  []Result `json:"results"`
}

func (cdo *CDO) collectResults() (chan *Result, error) {
	rChan := make(chan *Result, 10)

	go func() {
		for _, result := range cdo.Results {
			rChan <- &result
		}

		close(rChan)
	}()

	return rChan, nil
}

type Metadata struct {
	Resultset ResultSet `json:"resultset"`
}

type ResultSet struct {
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

func FetchDataFromStationForTimeSpan(station string, ts noaa.TimeSpan, token string) (chan *Result, error) {
	u, _ := url.Parse(BASE_URL + "/data")

	startdate := ts.Begin.Format("2006-01-02")
	enddate := ts.End.Format("2006-01-02")

	q := u.Query()
	q.Set("datasetid", "GHCND")
	q.Set("limit", "1000")
	q.Set("stationid", station)
	q.Set("startdate", startdate)
	q.Set("enddate", enddate)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("token", token)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var cdo CDO
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&cdo)

	if err != nil {
		return nil, err
	}

	rChan, err := cdo.collectResults()

	if err != nil {
		return nil, err
	}

	return rChan, nil
}
