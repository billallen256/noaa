package ndfd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFetchAndDecode(t *testing.T) {
	resp, err := http.Get("http://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php?whichClient=NDFDgen&lat=38.99&lon=-77.01&listLatLon=&lat1=&lon1=&lat2=&lon2=&resolutionSub=&listLat1=&listLon1=&listLat2=&listLon2=&resolutionList=&endPoint1Lat=&endPoint1Lon=&endPoint2Lat=&endPoint2Lon=&listEndPoint1Lat=&listEndPoint1Lon=&listEndPoint2Lat=&listEndPoint2Lon=&zipCodeList=&listZipCodeList=&centerPointLat=&centerPointLon=&distanceLat=&distanceLon=&resolutionSquare=&listCenterPointLat=&listCenterPointLon=&listDistanceLat=&listDistanceLon=&listResolutionSquare=&citiesLevel=&listCitiesLevel=&sector=&gmlListLatLon=&featureType=&requestedTime=&startTime=&endTime=&compType=&propertyName=&product=time-series&begin=2004-01-01T00%3A00%3A00&end=2018-07-27T00%3A00%3A00&Unit=e&maxt=maxt&mint=mint&temp=temp&qpf=qpf&pop12=pop12&snow=snow&dew=dew&wspd=wspd&wdir=wdir&sky=sky&wx=wx&waveh=waveh&icons=icons&rh=rh&appt=appt&incw34=incw34&incw50=incw50&incw64=incw64&cumw34=cumw34&cumw50=cumw50&cumw64=cumw64&critfireo=critfireo&dryfireo=dryfireo&conhazo=conhazo&ptornado=ptornado&phail=phail&ptstmwinds=ptstmwinds&pxtornado=pxtornado&pxhail=pxhail&pxtstmwinds=pxtstmwinds&ptotsvrtstm=ptotsvrtstm&pxtotsvrtstm=pxtotsvrtstm&tmpabv14d=tmpabv14d&tmpblw14d=tmpblw14d&tmpabv30d=tmpabv30d&tmpblw30d=tmpblw30d&tmpabv90d=tmpabv90d&tmpblw90d=tmpblw90d&prcpabv14d=prcpabv14d&prcpblw14d=prcpblw14d&prcpabv30d=prcpabv30d&prcpblw30d=prcpblw30d&prcpabv90d=prcpabv90d&prcpblw90d=prcpblw90d&precipa_r=precipa_r&sky_r=sky_r&td_r=td_r&temp_r=temp_r&wdir_r=wdir_r&wspd_r=wspd_r&wwa=wwa&wgust=wgust&iceaccum=iceaccum&maxrh=maxrh&minrh=minrh&Submit=Submit")

	if err != nil {
		t.Errorf("%s", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("%s", err)
	}

	resp.Body.Close()
	dwml, err := Unmarshal(body)

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(dwml)
}
