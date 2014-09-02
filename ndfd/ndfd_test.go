package ndfd

import (
	"fmt"
	"testing"
)

var ndfdGlobal NDFD

func TestFetchAndDecode(t *testing.T) {
	n, err := FetchNDFD(39.640102, -106.374332)

	if err != nil {
		t.Errorf("%s", err)
	}

	ndfdGlobal = n
	fmt.Printf("%+v\n", n.Dwml.Data.Parameters)
}

func TestDataLocation(t *testing.T) {
	fmt.Printf("%+v\n", ndfdGlobal.Dwml.Data.Location)
}

func TestHourlyVals(t *testing.T) {
	timeLayout, units, vals, err := ndfdGlobal.Dwml.Data.Parameters.HourlyTemperatures()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)

	timeLayout, units, vals, err = ndfdGlobal.Dwml.Data.Parameters.HourlyDewPoints()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)

	timeLayout, units, vals, err = ndfdGlobal.Dwml.Data.Parameters.HourlyCloudAmounts()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)

	timeLayout, units, vals, err = ndfdGlobal.Dwml.Data.Parameters.HourlyLiquidPrecip()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)

	timeLayout, units, vals, err = ndfdGlobal.Dwml.Data.Parameters.HourlyWindSpeeds()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)

	timeLayout, units, vals, err = ndfdGlobal.Dwml.Data.Parameters.HourlyWindDirections()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)

	timeLayout, units, vals, err = ndfdGlobal.Dwml.Data.Parameters.HourlySnowAmounts()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(timeLayout)
	fmt.Println(units)
	fmt.Println(vals)
}

func TestTimeSpanConditions(t *testing.T) {
	condChan, err := ndfdGlobal.Dwml.collectConditions()

	if err != nil {
		t.Errorf("%s", err)
	}

	for c := range condChan {
		fmt.Printf("%+v\n", c)
	}
}
