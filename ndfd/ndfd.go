package ndfd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type NDFD struct {
	SourceURL             string
	Dwml                  DWML
	TimeSpanCollectionMap map[string][]TimeSpan
}

func FetchNDFD() (NDFD, error) {
	sourceURL := "http://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php?whichClient=NDFDgen&lat=38.99&lon=-77.01&listLatLon=&lat1=&lon1=&lat2=&lon2=&resolutionSub=&listLat1=&listLon1=&listLat2=&listLon2=&resolutionList=&endPoint1Lat=&endPoint1Lon=&endPoint2Lat=&endPoint2Lon=&listEndPoint1Lat=&listEndPoint1Lon=&listEndPoint2Lat=&listEndPoint2Lon=&zipCodeList=&listZipCodeList=&centerPointLat=&centerPointLon=&distanceLat=&distanceLon=&resolutionSquare=&listCenterPointLat=&listCenterPointLon=&listDistanceLat=&listDistanceLon=&listResolutionSquare=&citiesLevel=&listCitiesLevel=&sector=&gmlListLatLon=&featureType=&requestedTime=&startTime=&endTime=&compType=&propertyName=&product=time-series&begin=2004-01-01T00%3A00%3A00&end=2018-07-27T00%3A00%3A00&Unit=e&maxt=maxt&mint=mint&temp=temp&qpf=qpf&pop12=pop12&snow=snow&dew=dew&wspd=wspd&wdir=wdir&sky=sky&wx=wx&waveh=waveh&icons=icons&rh=rh&appt=appt&incw34=incw34&incw50=incw50&incw64=incw64&cumw34=cumw34&cumw50=cumw50&cumw64=cumw64&critfireo=critfireo&dryfireo=dryfireo&conhazo=conhazo&ptornado=ptornado&phail=phail&ptstmwinds=ptstmwinds&pxtornado=pxtornado&pxhail=pxhail&pxtstmwinds=pxtstmwinds&ptotsvrtstm=ptotsvrtstm&pxtotsvrtstm=pxtotsvrtstm&tmpabv14d=tmpabv14d&tmpblw14d=tmpblw14d&tmpabv30d=tmpabv30d&tmpblw30d=tmpblw30d&tmpabv90d=tmpabv90d&tmpblw90d=tmpblw90d&prcpabv14d=prcpabv14d&prcpblw14d=prcpblw14d&prcpabv30d=prcpabv30d&prcpblw30d=prcpblw30d&prcpabv90d=prcpabv90d&prcpblw90d=prcpblw90d&precipa_r=precipa_r&sky_r=sky_r&td_r=td_r&temp_r=temp_r&wdir_r=wdir_r&wspd_r=wspd_r&wwa=wwa&wgust=wgust&iceaccum=iceaccum&maxrh=maxrh&minrh=minrh&Submit=Submit"

	resp, err := http.Get(sourceURL)

	if err != nil {
		return NDFD{}, err
	}

	if resp.StatusCode != 200 {
		return NDFD{}, errors.New(fmt.Sprintf("Received error %d from %s", resp.StatusCode, sourceURL))
	}

	var dwml DWML
	decoder := xml.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&dwml)

	if err != nil {
		return NDFD{}, err
	}

	tsMap, err := generateTimeSpanCollectionMap(&dwml)

	if err != nil {
		return NDFD{}, err
	}

	return NDFD{sourceURL, dwml, tsMap}, nil
}

func generateTimeSpanCollectionMap(dwml *DWML) (map[string][]TimeSpan, error) {
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
