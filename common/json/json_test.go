package json

import "testing"

func TestVerifyStringIsJSON(t *testing.T) {
	var objJSON = `{"a":1}`
	if !VerifyStringIsJSON(objJSON) {
		t.Fail()
	}
	var arrJSON = `[{"1": 2}]`
	if !VerifyStringIsJSON(arrJSON) {
		t.Fail()
	}
}
