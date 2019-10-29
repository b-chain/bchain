package main

import (
	"fmt"
	"testing"
)

func TestPurifyJS(t *testing.T) {
	tmp := []byte(`
/*
	multiple line comments  0
	multiple line comments  1
*/
var date=new Date();
var day=date.getDate();


date.setDate(day+200);
//
//single line comment
alert(date);
`)
	min, err := MinJS(tmp)
	if err != nil {
		fmt.Println("min js code is error", err)
	} else {
		fmt.Printf("last code is %s\n", min)
	}

	//fmt.Println(string(min))
}
