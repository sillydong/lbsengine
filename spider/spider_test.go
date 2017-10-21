package spider

import (
	"fmt"
	"testing"
)

func TestGetGaodeData(t *testing.T) {
	arrData := GetPOIData("1")
	fmt.Println("长度=", len(arrData), arrData)
	t.Log(arrData)
}
