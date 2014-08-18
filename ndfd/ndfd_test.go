package ndfd

import (
	"fmt"
	"testing"
)

var ndfdGlobal NDFD

func TestFetchAndDecode(t *testing.T) {
	n, err := FetchNDFD(39.0, -104.0)

	if err != nil {
		t.Errorf("%s", err)
	}

	ndfdGlobal = n
	fmt.Printf("%+v\n", n.Dwml.Data.Parameters)
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
	conditions, err := ndfdGlobal.Dwml.collectConditions()

	if err != nil {
		t.Errorf("%s", err)
	}

	for _, c := range conditions {
		fmt.Printf("%+v\n", c)
	}
}
