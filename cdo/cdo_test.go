package cdo

import (
	"fmt"
	"github.com/gershwinlabs/noaa"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFetchNewYork2014(t *testing.T) {
	station := "GHCND:USW00094728"
	token := strings.TrimSpace(os.Getenv("NOAA_TOKEN"))
	begin := time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC)
	ts := noaa.TimeSpan{begin, end}
	rChan, err := FetchDataFromStationForTimeSpan(station, ts, token)

	if err != nil {
		t.Errorf("%s", err)
	}

	for r := range rChan {
		fmt.Printf("%+v\n", r)
	}
}
