package time

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	date := time.Date(2022, time.April, 8, 13, 11, 13, 123456789, CSTZone)

	expectedList := map[string]string{
		Format:       "2022-04-08 13:11:13",
		RFC3339:      "2022-04-08T13:11:13+08:00",
		RFC3339Milli: "2022-04-08T13:11:13.123+08:00",
		RFC3339Micro: "2022-04-08T13:11:13.123456+08:00",
		RFC3339Nano:  "2022-04-08T13:11:13.123456789+08:00",
		Kitchen:      "1:11PM",
		KitchenSec:   "1:11:13PM",
		Date:         "2022-04-08",
		Date2:        "2022/04/08",
		Stamp:        "13:11:13",
		StampMilli:   "13:11:13.123",
		StampMicro:   "13:11:13.123456",
		StampNano:    "13:11:13.123456789",
	}
	for format, expected := range expectedList {
		actual := date.Format(format)
		if actual != expected {
			t.Errorf("time.Format(%v, %v) = %v, expected %v", t, format, actual, expected)
		}
	}
}

func TestMinMax(t *testing.T) {
	h1, _ := time.ParseDuration("1h")
	h2, _ := time.ParseDuration("2h")
	if h := Min(h1, h2); h != h1 {
		t.Errorf("Min(%v, %v) = %v, expected %v", h1, h2, h1, "1h")
	}
	if h := Max(h1, h2); h != h2 {
		t.Errorf("Max(%v, %v) = %v, expected %v", h1, h2, h2, "2h")
	}
}

func TestAdd(t *testing.T) {
	h1, _ := time.ParseDuration("1h")
	h2, _ := time.ParseDuration("2h")
	h3, _ := time.ParseDuration("3h")
	if h := Add(h1, h2); h != h3 {
		t.Errorf("Add(%v, %v) = %v, expected %v", h1, h2, h3, "2h")
	}

}
