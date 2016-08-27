package server
import (
	"testing"
	"net/url"
	"time"
)

func TestParams(t *testing.T) {
	params := Params{make(url.Values)}
	params.Set("stringkey", "a string")
	params.Set("intkey", "120")
	params.Set("timekey", "2010-01-03")
	params.Set("floatkey", "1.33")

	intVal, err := params.GetInt64("intkey")
	if err != nil || intVal != 120 {
		t.Fatal(err, intVal)
	}
	floatVal, err := params.GetFloat64("floatkey")
	if err != nil || floatVal != 1.33 {
		t.Fatal(err, intVal)
	}

	_, err = params.GetTime("timekey", "20060102")
	t.Log(err)
	if err == nil {
		t.Fatal("should fail to parse time")
	}
	timeVal, err := params.GetTime("timekey", "2006-01-02")
	if err != nil || timeVal != time.Date(2010, 1, 3, 0, 0, 0, 0, time.Local) {
		t.Fatal(err, timeVal)
	}
}
