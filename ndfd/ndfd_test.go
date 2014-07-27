package ndfd

import (
	"fmt"
	"testing"
)

var ndfdGlobal NDFD

func TestFetchAndDecode(t *testing.T) {
	n, err := FetchNDFD()

	if err != nil {
		t.Errorf("%s", err)
	}

	ndfdGlobal = n

	fmt.Println(n.Dwml)

}
