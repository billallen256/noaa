package cdo

import (
	"fmt"
	"github.com/gershwinlabs/noaa"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSubTimeSpans(t *testing.T) {
	begin := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	overall := noaa.TimeSpan{begin, end}
	timeSpans := subTimeSpans(overall)

	if len(timeSpans) != 4 {
		t.Errorf("%d TimeSpans returned, but should have received 4 TimeSpans\n", len(timeSpans))
	}

	if timeSpans[0].Begin != begin {
		t.Errorf("First TimeSpan has incorrect end %v", timeSpans[0].Begin)
	}

	if timeSpans[len(timeSpans)-1].End != end {
		t.Errorf("Last TimeSpan has incorrect end %v", timeSpans[len(timeSpans)-1].End)
	}
}

func TestFetchNewYork2014(t *testing.T) {
	station := "GHCND:USW00094728"
	token := strings.TrimSpace(os.Getenv("NOAA_TOKEN"))
	begin := time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC)
	ts := noaa.TimeSpan{begin, end}
	rChan, err := FetchDataFromStationForTimeSpan(station, ts, token)
	numResultsFetched := 0

	if err != nil {
		t.Errorf("%s", err)
	}

	for r := range rChan {
		numResultsFetched++
		fmt.Printf("%+v\n", r)
	}

	if numResultsFetched == 0 {
		t.Errorf("No results fetched")
	}
}

func TestFetchNewYorkOverLimit(t *testing.T) {
	station := "GHCND:USW00094728"
	token := strings.TrimSpace(os.Getenv("NOAA_TOKEN"))
	begin := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2013, 12, 31, 0, 0, 0, 0, time.UTC)
	ts := noaa.TimeSpan{begin, end}
	rChan, err := FetchDataFromStationForTimeSpan(station, ts, token)
	numResultsFetched := 0

	if err != nil {
		t.Errorf("%s", err)
	}

	for r := range rChan {
		numResultsFetched++
		fmt.Printf("%+v\n", r)
	}

	if numResultsFetched == 0 {
		t.Errorf("No results fetched")
	}
}
