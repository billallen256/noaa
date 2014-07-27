package ndfd

import (
	"encoding/xml"
	"time"
)

type TimeSpan struct {
	Begin time.Time
	End   time.Time
}

type DWML struct {
	Head Head `xml:"head"`
	Data Data `xml:"data"`
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
	ApplicableLocation         string                           `xml:"applicable-location,attr"`
	Temperature                []DataParametersSection          `xml:"temperature"`
	Precipitation              []DataParametersSection          `xml:"precipitation"`
	WindSpeed                  []DataParametersSection          `xml:"wind-speed"`
	Direction                  []DataParametersSection          `xml:"direction"`
	CloudAmount                []DataParametersSection          `xml:"cloud-amount"`
	ProbabilityOfPrecipitation []DataParametersSection          `xml:"probability-of-precipitation"`
	FireWeather                []DataParametersSection          `xml:"fire-weather"`
	ConvectiveHazard           []DataParametersConvectiveHazard `xml:"convective-hazard"`
	ClimateAnomaly             []DataParametersClimateAnomaly   `xml:"climate-anomoly"`
	Humidity                   []DataParametersSection          `xml:"humidity"`
	Weather                    DataParametersWeather            `xml:"weather"`
	ConditionsIcon             DataParametersConditionsIcon     `xml:"conditions-icon"`
	Hazards                    DataParametersHazards            `xml:"hazards"`
	WaterState                 DataParametersWaterState         `xml:"water-state"`
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

func generateTimeSpanCollectionMap(dwml DWML) (map[string][]TimeSpan, error) {
	m := make(map[string][]TimeSpan)

	for _, timeLayout := range dwml.Data.TimeLayouts {
		numStartTimes := len(timeLayout.StartValidTimes)
		numEndTimes := len(timeLayout.EndValidTimes)
		arr := make([]TimeSpan, numStartTimes)

		for i := 0; i < numStartTimes; i++ {
			begin, err := time.Parse(time.RFC3339, timeLayout.StartValidTimes[i])

			if err != nil {
				return m, err
			}

			end := begin

			if numEndTimes == numStartTimes {
				end, err = time.Parse(time.RFC3339, timeLayout.EndValidTimes[i])

				if err != nil {
					return m, err
				}
			}

			arr[i] = TimeSpan{begin, end}
		}

		m[timeLayout.LayoutKey] = arr
	}

	return m, nil
}

func Unmarshal(body []byte) (DWML, error) {
	var dwml DWML

	err := xml.Unmarshal(body, &dwml)

	if err != nil {
		return DWML{}, err
	}

	return dwml, nil
}
