package jsre

//js contract must implemente these standard api
var getName = "getName()"
var getVersion = "getVersion()"

var endDot = ";"

func GetJsStandardApis() string {
	apis := getName + endDot
	apis += getVersion + endDot
	return apis

}
