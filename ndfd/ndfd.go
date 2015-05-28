package ndfd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gershwinlabs/noaa"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	sourceURLFormat = "http://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php?whichClient=NDFDgen&lat=%f&lon=%f&product=time-series&begin=%s&end=%s&Unit=m&maxt=maxt&mint=mint&temp=temp&qpf=qpf&pop12=pop12&snow=snow&dew=dew&wspd=wspd&wdir=wdir&sky=sky&wx=wx&waveh=waveh&icons=icons&rh=rh&appt=appt&incw34=incw34&incw50=incw50&incw64=incw64&cumw34=cumw34&cumw50=cumw50&cumw64=cumw64&critfireo=critfireo&dryfireo=dryfireo&conhazo=conhazo&ptornado=ptornado&phail=phail&ptstmwinds=ptstmwinds&pxtornado=pxtornado&pxhail=pxhail&pxtstmwinds=pxtstmwinds&ptotsvrtstm=ptotsvrtstm&pxtotsvrtstm=pxtotsvrtstm&tmpabv14d=tmpabv14d&tmpblw14d=tmpblw14d&tmpabv30d=tmpabv30d&tmpblw30d=tmpblw30d&tmpabv90d=tmpabv90d&tmpblw90d=tmpblw90d&prcpabv14d=prcpabv14d&prcpblw14d=prcpblw14d&prcpabv30d=prcpabv30d&prcpblw30d=prcpblw30d&prcpabv90d=prcpabv90d&prcpblw90d=prcpblw90d&precipa_r=precipa_r&sky_r=sky_r&td_r=td_r&temp_r=temp_r&wdir_r=wdir_r&wspd_r=wspd_r&wwa=wwa&wgust=wgust&iceaccum=iceaccum&maxrh=maxrh&minrh=minrh&Submit=Submit"
)

type NDFD struct {
	SourceURL  string
	Dwml       *DWML
	Conditions chan Condition
}

type Condition struct {
	Name  string
	Value float64
	Units string
	Hour  time.Time
	Lat   float64
	Lon   float64
}

func FetchNDFD(lat, lon float64) (NDFD, error) {
	return FetchNDFDWithClient(http.DefaultClient, lat, lon)
}

func FetchNDFDWithClient(client *http.Client, lat, lon float64) (NDFD, error) {
	b := time.Now().UTC().Add(time.Duration(-10*24) * time.Hour)
	e := time.Now().UTC().Add(time.Duration(10*24) * time.Hour)
	return FetchNDFDWithClientForTimeSpan(client, noaa.TimeSpan{b, e}, lat, lon)
}

func FetchNDFDCurrent(lat, lon float64) (NDFD, error) {
	return FetchNDFDCurrentWithClient(http.DefaultClient, lat, lon)
}

func FetchNDFDCurrentWithClient(client *http.Client, lat, lon float64) (NDFD, error) {
	b := time.Now().UTC().Add(time.Duration(-1) * time.Hour)
	e := time.Now().UTC().Add(time.Duration(1) * time.Hour)
	return FetchNDFDWithClientForTimeSpan(client, noaa.TimeSpan{b, e}, lat, lon)
}

func FetchNDFDForecast(lat, lon float64) (NDFD, error) {
	return FetchNDFDForecastWithClient(http.DefaultClient, lat, lon)
}

func FetchNDFDForecastWithClient(client *http.Client, lat, lon float64) (NDFD, error) {
	b := time.Now().UTC()
	e := time.Now().UTC().Add(time.Duration(7*24) * time.Hour)
	return FetchNDFDWithClientForTimeSpan(client, noaa.TimeSpan{b, e}, lat, lon)
}

func FetchNDFDWithClientForTimeSpan(client *http.Client, ts noaa.TimeSpan, lat, lon float64) (NDFD, error) {
	b := url.QueryEscape(ts.Begin.Format("2006-01-02T15:04:05"))
	e := url.QueryEscape(ts.End.Format("2006-01-02T15:04:05"))
	sourceURL := fmt.Sprintf(sourceURLFormat, lat, lon, b, e)
	resp, err := client.Get(sourceURL)

	if err != nil {
		return NDFD{}, err
	}

	return processNDFDResponse(resp, sourceURL)
}

func processNDFDResponse(resp *http.Response, sourceURL string) (NDFD, error) {
	if resp.StatusCode != 200 {
		return NDFD{}, errors.New(fmt.Sprintf("Received error %d from %s", resp.StatusCode, sourceURL))
	}

	var dwml DWML
	decoder := xml.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err := decoder.Decode(&dwml)

	if err != nil {
		return NDFD{}, err
	}

	condChan, err := dwml.collectConditions()

	if err != nil {
		return NDFD{}, err
	}

	return NDFD{sourceURL, &dwml, condChan}, nil
}

type DWML struct {
	Head Head `xml:"head"`
	Data Data `xml:"data"`
}

func (dwml *DWML) generateTimeSpanLayoutMap() (map[string][]noaa.TimeSpan, error) {
	m := make(map[string][]noaa.TimeSpan)

	for _, timeLayout := range dwml.Data.TimeLayouts {
		numStartTimes := len(timeLayout.StartValidTimes)
		numEndTimes := len(timeLayout.EndValidTimes)
		arr := make([]noaa.TimeSpan, numStartTimes)

		for i := 0; i < numStartTimes; i++ {
			begin, err := time.ParseInLocation(time.RFC3339, timeLayout.StartValidTimes[i], time.UTC)

			if err != nil {
				return m, err
			}

			end := begin

			if numEndTimes == numStartTimes {
				end, err = time.ParseInLocation(time.RFC3339, timeLayout.EndValidTimes[i], time.UTC)

				if err != nil {
					return m, err
				}
			}

			arr[i] = noaa.TimeSpan{begin, end}
		}

		m[timeLayout.LayoutKey] = arr
	}

	return m, nil
}

func (dwml *DWML) collectConditions() (chan Condition, error) {
	tsMap, err := dwml.generateTimeSpanLayoutMap()
	condChan := make(chan Condition, 10)

	if err != nil {
		return condChan, err
	}

	go func() {
		lat := dwml.Data.Location.Point.Latitude
		lon := dwml.Data.Location.Point.Longitude

		layout, units, vals, err := dwml.Data.Parameters.HourlyTemperatures()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"temp", val, units, hour, lat, lon}
				}
			}
		}

		layout, units, vals, err = dwml.Data.Parameters.HourlyDewPoints()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"dewpoint", val, units, hour, lat, lon}
				}
			}
		}

		layout, units, vals, err = dwml.Data.Parameters.HourlyCloudAmounts()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"clouds", val, units, hour, lat, lon}
				}
			}
		}

		layout, units, vals, err = dwml.Data.Parameters.HourlyLiquidPrecip()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"precip", val, units, hour, lat, lon}
				}
			}
		}

		layout, units, vals, err = dwml.Data.Parameters.HourlyWindSpeeds()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"windspeed", val, units, hour, lat, lon}
				}
			}
		}

		layout, units, vals, err = dwml.Data.Parameters.HourlyWindDirections()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"winddir", val, units, hour, lat, lon}
				}
			}
		}

		layout, units, vals, err = dwml.Data.Parameters.HourlySnowAmounts()

		if err == nil {
			for i, val := range vals {
				if math.IsNaN(val) {
					continue
				}

				for _, hour := range tsMap[layout][i].Hours() {
					condChan <- Condition{"snow", val, units, hour, lat, lon}
				}
			}
		}

		close(condChan)
	}()

	return condChan, nil
}

type Head struct {
	Product HeadProduct `xml:"product"`
	Source  HeadSource  `xml:"source"`
}

type HeadProduct struct {
	SrsName         string                  `xml:"srsName,attr"`
	ConciseName     string                  `xml:"concise-name,attr"`
	OperationalMode string                  `xml:"operational-mode,attr"`
	Title           string                  `xml:"title"`
	Field           string                  `xml:"field"`
	Category        string                  `xml:"category"`
	CreationDate    HeadProductCreationDate `xml:"creation-date"`
}

type HeadProductCreationDate struct {
	RefreshFreq string `xml:"refresh-frequency,attr"`
	Value       string `xml:",chardata"`
}

type HeadSource struct {
	MoreInformation  string                     `xml:"more-information"`
	ProductionCenter HeadSourceProductionCenter `xml:"production-center"`
	Disclaimer       string                     `xml:"disclaimer"`
	Credit           string                     `xml:"credit"`
	CreditLogo       string                     `xml:"credit-logo"`
	Feedback         string                     `xml:"feedback"`
}

type HeadSourceProductionCenter struct {
	Value     string `xml:",chardata"`
	SubCenter string `xml:"sub-center"`
}

type Data struct {
	Location               DataLocation               `xml:"location"`
	MoreWeatherInformation DataMoreWeatherInformation `xml:"moreWeatherInformation"`
	TimeLayouts            []DataTimeLayout           `xml:"time-layout"`
	Parameters             DataParameters             `xml:"parameters"`
}

type DataLocation struct {
	LocationKey string            `xml:"location-key"`
	Point       DataLocationPoint `xml:"point"`
}

type DataLocationPoint struct {
	Latitude  float64 `xml:"latitude,attr"`
	Longitude float64 `xml:"longitude,attr"`
}

type DataMoreWeatherInformation struct {
	ApplicableLocation string `xml:"applicable-location,attr"`
	Value              string `xml:",chardata"`
}

type DataTimeLayout struct {
	TimeCoordinate  string   `xml:"time-coordinate,attr"`
	Summarization   string   `xml:"summarization,attr"`
	LayoutKey       string   `xml:"layout-key"`
	StartValidTimes []string `xml:"start-valid-time"`
	EndValidTimes   []string `xml:"end-valid-time"`
}

type DataParameters struct {
	ApplicableLocation           string                           `xml:"applicable-location,attr"`
	Temperatures                 []DataParametersSection          `xml:"temperature"`
	Precipitations               []DataParametersSection          `xml:"precipitation"`
	WindSpeeds                   []DataParametersSection          `xml:"wind-speed"`
	Directions                   []DataParametersSection          `xml:"direction"`
	CloudAmounts                 []DataParametersSection          `xml:"cloud-amount"`
	ProbabilitiesOfPrecipitation []DataParametersSection          `xml:"probability-of-precipitation"`
	FireWeathers                 []DataParametersSection          `xml:"fire-weather"`
	ConvectiveHazards            []DataParametersConvectiveHazard `xml:"convective-hazard"`
	ClimateAnomalies             []DataParametersClimateAnomaly   `xml:"climate-anomoly"`
	Humidities                   []DataParametersSection          `xml:"humidity"`
	Weathers                     DataParametersWeather            `xml:"weather"`
	ConditionsIcon               DataParametersConditionsIcon     `xml:"conditions-icon"`
	Hazards                      DataParametersHazards            `xml:"hazards"`
	WaterState                   DataParametersWaterState         `xml:"water-state"`
}

func GetParametersSection(dps []DataParametersSection, name string) (string, string, []float64, error) {
	timeLayout := "unknown"
	units := "unknown"
	vals := make([]float64, 0, 64)

	for _, t := range dps {
		if t.Type == name {
			timeLayout = t.TimeLayout
			units = t.Units

			for _, v := range t.Values {
				i, err := strconv.ParseFloat(v, 64)

				if err != nil {
					vals = append(vals, math.NaN())
				} else {
					vals = append(vals, i)
				}
			}
		}
	}

	if timeLayout == "unknown" {
		return timeLayout, units, []float64{}, errors.New(fmt.Sprintf("Could not find section %s", name))
	}

	return timeLayout, units, vals, nil
}

func (dp DataParameters) HourlyTemperatures() (string, string, []float64, error) {
	return GetParametersSection(dp.Temperatures, "hourly")
}

func (dp DataParameters) HourlyDewPoints() (string, string, []float64, error) {
	return GetParametersSection(dp.Temperatures, "dew point")
}

func (dp DataParameters) HourlyCloudAmounts() (string, string, []float64, error) {
	return GetParametersSection(dp.CloudAmounts, "total")
}

func (dp DataParameters) HourlyLiquidPrecip() (string, string, []float64, error) {
	return GetParametersSection(dp.Precipitations, "liquid")
}

func (dp DataParameters) HourlyWindSpeeds() (string, string, []float64, error) {
	return GetParametersSection(dp.WindSpeeds, "sustained")
}

func (dp DataParameters) HourlyWindDirections() (string, string, []float64, error) {
	return GetParametersSection(dp.Directions, "wind")
}

func (dp DataParameters) HourlySnowAmounts() (string, string, []float64, error) {
	return GetParametersSection(dp.Precipitations, "snow")
}

type DataParametersSection struct {
	Type       string   `xml:"type,attr"`
	Units      string   `xml:"units,attr"`
	TimeLayout string   `xml:"time-layout,attr"`
	Name       string   `xml:"name"`
	Values     []string `xml:"value"`
}

type DataParametersConvectiveHazard struct {
	Outlook         DataParametersSection `xml:"outlook"`
	SevereComponent DataParametersSection `xml:"severe-component"`
}

type DataParametersClimateAnomaly struct {
	Weekly   DataParametersSection `xml:"weekly"`
	Monthly  DataParametersSection `xml:"monthly"`
	Seasonal DataParametersSection `xml:"seasonal"`
}

type DataParametersWeather struct {
	TimeLayout        string                            `xml:"time-layout,attr"`
	Name              string                            `xml:"name"`
	WeatherConditions []DataParametersWeatherConditions `xml:"weather-conditions"`
}

type DataParametersWeatherConditions struct {
	Value DataParametersWeatherConditionsValue `xml:"value"`
}

type DataParametersWeatherConditionsValue struct {
	Coverage    string                                         `xml:"coverage,attr"`
	Intensity   string                                         `xml:"intensity,attr"`
	WeatherType string                                         `xml:"weather-type,attr"`
	Qualifier   string                                         `xml:qualifier,attr"`
	Visibility  DataParametersWeatherConditionsValueVisibility `xml:"visibility"`
}

type DataParametersWeatherConditionsValueVisibility struct {
	Units string `xml:"units,attr"`
	Value string `xml:",chardata"`
}

type DataParametersConditionsIcon struct {
	Type       string   `xml:"type,attr"`
	TimeLayout string   `xml:"time-layout,attr"`
	Name       string   `xml:"name"`
	IconLink   []string `xml"icon-link"`
}

type DataParametersHazards struct {
	TimeLayout       string `xml:"time-layout,attr"`
	Name             string `xml:"name"`
	HazardConditions []interface{}
}

type DataParametersWaterState struct {
	TimeLayout string `xml:"time-layout,attr"`
	Waves      DataParametersSection
}
