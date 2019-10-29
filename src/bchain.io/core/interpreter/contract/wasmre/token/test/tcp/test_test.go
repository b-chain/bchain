package tcp

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type updateCar struct {
	ParkId string `json:"parkId"`
	ApiKey string `json:"apiKey"`

	UserId    string `json:"userId"`
	UserRoom  string `json:"userRoom"`
	UserName  string `json:"userName"`
	UserAddr  string `json:"userAddr"`
	UserPhone string `json:"userPhone"`

	Cars       []Cars      `json:"cars"`
	ParkSpaces []ParkSpace `json:"parkSpaces"`
	Amount     float64     `json:"amount"`
}

type ParkSpace struct {
	Area      string `json:"area"`
	ParkSpace string `json:"parkSpace"`
}

type getCar struct {
	ParkId string `json:"parkId"`
	ApiKey string `json:"apiKey"`
	UserId string `json:"userId"`
}

type getCarRlt struct {
	UserId    string  `json:"userId"`
	UserRoom  string  `json:"userRoom"`
	UserName  string  `json:"userName"`
	UserAddr  string  `json:"userAddr"`
	UserPhone string  `json:"userPhone"`
	Amount    float64 `json:"amount"`

	Cars       []Cars   `json:"cars"`
	ParkSpaces []string `json:"parkSpace"`
}

type Cars struct {
	Plate        string `json:"plate"`
	Area         string `json:"area"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	Type         int    `json:"type"`
	ChargeRuleID string `json:"chargeRuleId"`
}
type GapFee [24]float64
type standChargeRule struct {
	StartHour   int  `json:"startHour"`
	StartMinute int     `json:"startMinute"`
	StartFee    float64 `json:"startFee"`
	GapMinute   int     `json:"gapMinute"`
	GapFeeRate  float64 `json:"gapFeeRate"`
}
type ChargeRule struct {
	Type       int             `json:"type"` //0，free, 1, timeLenTable, 2, timeTable
	FreeMinute int             `json:"freeMinute"`
	GapFee     *GapFee          `json:"gapFee"`
	FeeMax     float64         `json:"feeMax"`
	DayCr      *standChargeRule `json:"dayCr"`
	NightCr    *standChargeRule `json:"nightCr"`
}

func TestApi(t *testing.T) {
	//fg:= GapFee{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24}
	//cr := &ChargeRule{1, 15, fg, 100}
	tm := time.Now()
	tmS := tm.Format("2006-1-2 15:04:05")
	tmE := tm.AddDate(1, 0, 1).Format("2006-1-2 15:04:05")
	car1 := Cars{"川A51234", "1号地库", tmS, tmE, 0, "vipa"}
	car2 := Cars{"川A51235", "1号地库", tmS, tmE, 0, "vipa"}
	car3 := Cars{"川A51236", "1号地库", tmS, tmE, 1, "vipa"}
	updatecar := updateCar{
		ParkId:    "42342342adf",
		ApiKey:    "asexeolklsrj",
		UserId:    "1-4-806",
		UserRoom:  "806",
		UserName:  "张三",
		UserAddr:  "南京东路11号",
		UserPhone: "13889899999",
		Amount:    1000,
	}
	updatecar.Cars = append(updatecar.Cars, car1, car2, car3)
	pa1 := ParkSpace{"1号地库", "s-00001"}
	pa2 := ParkSpace{"1号地库", "s-00002"}
	updatecar.ParkSpaces = append(updatecar.ParkSpaces, pa1, pa2)
	paraBytes, _ := json.MarshalIndent(&updatecar, "", "    ")
	fmt.Println(string(paraBytes))

	yy := &updateCar{}
	aaa := json.Unmarshal(paraBytes, yy)
	fmt.Println(aaa, yy)

	dd := make(map[string]interface{})
	dd["success"] = true
	dd["msg"] = "ok!"
	dd["code"] = 0
	dd["data"] = nil
	xx, _ := json.MarshalIndent(dd, "", "    ")
	fmt.Println(string(xx))

}

type areas struct {
	Id         string      `json:"id"`
	ChargeRuleId string `json:"chargeRuleId"`
	ChargeRuleIdYellow string `json:"chargeRuleIdYellow"`
}

type updateAreaChargeRule struct {
	ParkId     string `json:"parkId"`
	ApiKey     string `json:"apiKey"`
	Id         string `json:"id"`
	ChargeRuleId string  `json:"ChargeRuleId"`
	ChargeRuleIdYellow string `json:"chargeRuleIdYellow"`
}

type updateChargeRule struct {
	ParkId     string `json:"parkId"`
	ApiKey     string `json:"apiKey"`
	Id string  `json:"id"`
	ChargeRule *ChargeRule `json:"chargeRule"`
}

func TestApixx(t *testing.T) {
	dd := make(map[string]interface{})
	dd["success"] = true
	dd["msg"] = "ok!"
	dd["code"] = 0
	dd["data"] = nil

	//dcr := &standChargeRule{8,120, 5,60,2}
	//lcr := &standChargeRule{20,120, 5,60,1}
	//fg := &GapFee{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	//cr := &ChargeRule{1, 15, fg, 100,nil,nil}
	//cr2:= &ChargeRule{2, 15, nil, 100,dcr,lcr}
	aa := areas{
		Id:         "1号地库",
		ChargeRuleId: "tmpA",
	}
	bb := areas{
		Id:         "2号地库",
		ChargeRuleId: "tmpB",
	}
	as := make([]areas, 0)
	as = append(as, aa, bb)
	dd["data"] = as
	xx, _ := json.MarshalIndent(dd, "", "    ")
	fmt.Println(string(xx))

	uu := updateAreaChargeRule{"42342342adf", "asexeolklsrj", "1号地库", "tmpA","tmpB"}

	xx, _ = json.MarshalIndent(&uu, "", "    ")
	fmt.Println(string(xx))
}

func TestApixyx(t *testing.T) {
	dd := make(map[string]interface{})
	dd["success"] = true
	dd["msg"] = "ok!"
	dd["code"] = 0
	dd["data"] = nil

	dcr := &standChargeRule{8,120, 5,60,2}
	lcr := &standChargeRule{20,120, 5,60,1}
	fg := &GapFee{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	fg1 := &GapFee{5, 5, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50}
	cr := &ChargeRule{1, 15, fg1, 100,nil,nil}
	cr2:= &ChargeRule{2, 15, nil, 100,dcr,lcr}
	cr3:= &ChargeRule{1, 15, fg, 100,nil,nil}

	as := make(map[string]*ChargeRule, 0)
	as["tmpA"] = cr
	as["tmpB"] = cr2
	as["vipA"] = cr3
	dd["data"] = as
	xx, _ := json.MarshalIndent(dd, "", "    ")
	fmt.Println(string(xx))

	ee := make(map[string]interface{})
	err := json.Unmarshal(xx,&ee)
	fmt.Println(err,ee)

	uu := updateChargeRule{"42342342adf", "asexeolklsrj", "vipB", cr}

	xx, _ = json.MarshalIndent(&uu, "", "    ")
	fmt.Println(string(xx))
}

type CarFee struct {
	Area      string  `json:"area"`
	Plate     string  `json:"plate"`
	EnterTime string  `json:"enterTime"`
	CurTime   string  `json:"curTime"`
	TimeLen   int     `json:"timeLen"`
	Fee       float64 `json:"fee"`
	Memo      string  `json:"memo"`
}

func TestApixxx(t *testing.T) {
	dd := make(map[string]interface{})
	dd["success"] = true
	dd["msg"] = "ok!"
	dd["code"] = "order20190302111250川A12345"
	dd["data"] = nil

	tm := time.Now()
	tmS := tm.Format("2006-1-2 15:04:05")
	tmE := tm.Add(60 * time.Second).Format("2006-1-2 15:04:05")
	cf := &CarFee{
		EnterTime: tmS,
		CurTime:   tmE,
		TimeLen:   60,
		Fee:       64,
	}
	dd["data"] = cf
	xx, _ := json.MarshalIndent(dd, "", "    ")
	fmt.Println(string(xx))

}

func TestChan(t *testing.T) {
	ch := make(chan int, 0)
	close(ch)
	close(ch)
	//ch <- 1
	xx := <- ch
	fmt.Println(xx)
}
