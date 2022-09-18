package toolkit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var jsonTests = []struct {
	name          string
	json          string
	errorExpected bool
	maxSize       int
	allowUnknown  bool
}{
	{name: "good json", json: `{"foo":"bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: false},
	{name: "badly formed json", json: `{"foo"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "incorrect type", json: `{"foo":1}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "two json values", json: `{"foo":"bar"}{"bizz":"bazz"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "empty body", json: ``, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "json syntax error", json: `{"foo":"bar"`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "json unknown field", json: `{"fooo":"bazz"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "allow unknown fields", json: `{"bizz":"bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: true},
	{name: "missing field", json: `{"bizz":"bazz"}`, errorExpected: true, maxSize: 1024, allowUnknown: true},
	{name: "json too large", json: `{"foo":"bar"}`, errorExpected: true, maxSize: 4, allowUnknown: true},
	{name: "no json", json: `Foo bar`, errorExpected: true, maxSize: 1024, allowUnknown: true},
}

func TestToolkit_ReadJSON(t *testing.T) {
	var tools Tools

	for _, val := range jsonTests {
		tools.MaxJSONSize = val.maxSize
		tools.AllowUnknownFields = val.allowUnknown

		var testJSON struct {
			Foo string `json:"foo"`
		}

		request, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(val.json)))
		if err != nil {
			t.Log("Error:", err)
		}

		recorder := httptest.NewRecorder()

		err = tools.ReadJSON(recorder, request, &testJSON)

		if val.errorExpected && err == nil {
			t.Errorf("%s: error expected, none received", val.name)
		}

		if !val.errorExpected && err != nil {
			t.Errorf("%s: no error expected, error received: %s", val.name, err.Error())
		}

		request.Body.Close()
	}
}
