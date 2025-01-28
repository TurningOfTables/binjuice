package bin

import (
	"strconv"
	"testing"
	"time"
)

func TestNextCollectionData(t *testing.T) {
	res := NextCollectionData()
	if len(res) != 4 {
		t.Errorf("NextCollectionData() length (%v)", res)
	}
}

func TestNextCollectionDate(t *testing.T) {

	var tests = []struct {
		tYear  int
		tMonth time.Month
		tDay   int
		eYear  int
		eMonth time.Month
		eDay   int
	}{
		{2023, 7, 24, 2023, 7, 27},
		{2023, 7, 27, 2023, 7, 27},
		{2023, 7, 30, 2023, 8, 3},
		{2023, 9, 15, 2023, 9, 21},
		{2024, 4, 5, 2024, 4, 11},
		{2022, 2, 20, 2022, 2, 24},
	}

	for _, test := range tests {
		testDate := time.Date(test.tYear, test.tMonth, test.tDay, 1, 0, 0, 0, time.Local)
		res := NextCollectionDate(testDate)
		if res.Day() != test.eDay {
			t.Errorf("NextCollectionDate(%v) = %v", testDate, res)
		}

	}
}

func TestCurrentWeekYear(t *testing.T) {
	year, week := NextCollectionWeekYear()

	if len(strconv.Itoa(year)) != 4 {
		t.Errorf("CurrentWeekYear year length (%v)", year)
	}
	if week < 1 || week > 53 {
		t.Errorf("CurrentWeekYear week out of range (%v)", week)
	}
}

func TestWeeksInYear(t *testing.T) {
	var tests = []struct {
		input int
		want  int
	}{
		{2000, 53},
		{2004, 53},
		{2020, 53},
		{2032, 53},
		{2044, 53},
		{2001, 52},
		{2005, 52},
		{2011, 52},
		{2021, 52},
		{2023, 52},
	}

	for _, test := range tests {
		if got := WeeksInYear(test.input); got != test.want {
			t.Errorf("weeksInYear(%v) = %v", test.input, got)
		}
	}
}

func TestWeeksSinceAnchor(t *testing.T) {
	var tests = []struct {
		cYear int
		cWeek int
		aYear int
		aWeek int
		want  int
	}{
		{2023, 29, 2022, 35, 46},
		{2023, 29, 2022, 34, 47},
		{2023, 29, 2022, 16, 65},
		{2025, 4, 2018, 30, 340},
	}

	for _, test := range tests {
		got, err := WeeksSinceAnchor(test.cYear, test.cWeek, test.aYear, test.aWeek)
		if err != nil {
			t.Errorf("weeksSinceAnchor(%v, %v, %v, %v) = %v", test.cYear, test.cWeek, test.aYear, test.aWeek, err)
		}

		if got != test.want {
			t.Errorf("weeksSinceAnchor(%v, %v, %v, %v) = %v", test.cYear, test.cWeek, test.aYear, test.aWeek, got)
		}
	}
}

func TestIsBinCollected(t *testing.T) {
	var tests = []struct {
		bd   binData
		ws   int
		want bool
	}{
		{binData{Freq: 3}, 21, true},
		{binData{Freq: 2}, 48, true},
		{binData{Freq: 1}, 5, true},
		{binData{Freq: 3}, 22, false},
		{binData{Freq: 2}, 39, false},
	}

	for _, test := range tests {
		if res := IsBinCollected(test.bd, test.ws); res != test.want {
			t.Errorf("isBinCollected(%v, %v) = %v", test.bd, test.ws, test.want)
		}
	}
}
