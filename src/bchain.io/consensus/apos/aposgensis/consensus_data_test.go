package aposgensis

import (
	"testing"
	"bchain.io/common/types"
	"encoding/json"
	"fmt"
)

func TestXX(t *testing.T)  {
	xx := WeightInfos{}
	a := WeightInfo{types.HexToAddress("8ad432ee1ac9c75bf696620ad237c4484514347f"),100}
	xx = append(xx, a)
	jsonData,_:= json.MarshalIndent(&xx,"","	")
	fmt.Println(string(jsonData))
}
