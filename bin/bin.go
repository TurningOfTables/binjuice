package bin

import (
	"errors"
	"log"
	"time"
)

const (
	COLLECTIONDAY   = time.Thursday
	WEEKSINYEAR     = 52
	WEEKSINLEAPYEAR = 53
	ANCHORYEAR      = 2025
	ANCHORMONTH     = 1
	ANCHORDAY       = 16
)

type binData struct {
	Name         string
	FriendlyName string
	ImagePath    string
	Freq         int // bin is collected every <Freq> weeks
	Collected    bool
}

var bd = []binData{
	{Name: "Food", FriendlyName: "Useless Food", ImagePath: "food_bin.png", Freq: 1},
	{Name: "Recycling", FriendlyName: "Recycling", ImagePath: "recycling_bin.png", Freq: 2},
	{Name: "Grey", FriendlyName: "Waste", ImagePath: "grey_bin.png", Freq: 3},
	{Name: "Green", FriendlyName: "Garden", ImagePath: "green_bin.png", Freq: 2},
}

// NextCollectionData iterates through entries in a []binData and appends a Collected: bool to the bin.
//
// It returns the updated []binData
func NextCollectionData() []binData {
	aYear, aWeek, err := AnchorWeekYear()
	if err != nil {
		log.Fatalf("Error getting the week and year of the anchor date (%v)", err)
	}

	ncYear, ncWeek := NextCollectionWeekYear()

	wsAnchor, err := WeeksSinceAnchor(ncYear, ncWeek, aYear, aWeek)
	if err != nil {
		log.Fatal("error getting the weeks since the anchor date")
	}

	for k, bin := range bd {
		bd[k].Collected = IsBinCollected(bin, wsAnchor)
	}

	return bd
}

// NextCollectionDate uses the const COLLECTIONDAY which is a time.Weekday.
//
// It returns a time.Time of the next instance of that weekday.
func NextCollectionDate(t time.Time) time.Time {
	var ncd time.Time

	// It's already the collection day so return date
	if t.Weekday() == COLLECTIONDAY {
		return t
	}

	// Get nearest collection day
	ncd = t.AddDate(0, 0, int(COLLECTIONDAY)-int(t.Weekday()))

	// If nearest collection day is in the past, add a week
	if ncd.Before(t) {
		ncd = ncd.AddDate(0, 0, 7)
	}

	return ncd
}

// NextCollectionWeekYear returns the current week and year, or the following week and year if it's after COLLECTIONDAY
func NextCollectionWeekYear() (year int, week int) {
	t := time.Now()
	if time.Now().Weekday() > COLLECTIONDAY || time.Now().Weekday() == 0 {
		t = t.AddDate(0, 0, 7)
	}
	year, week = t.ISOWeek()
	return year, week
}

// AnchorWeekYear returns the week and year of the date specified by the constants.
// ANCHORYEAR, ANCHORMONTH, ANCHORDAY.
//
// It returns 0, 0, error if any constants are outside valid ranges.
func AnchorWeekYear() (awYear int, awWeek int, err error) {
	if ANCHORYEAR < 1900 || ANCHORYEAR > 9999 {
		return 0, 0, errors.New("value of ANCHORYEAR is outside range 1900-9999")
	}
	if ANCHORMONTH < 1 || ANCHORMONTH > 12 {
		return 0, 0, errors.New("value of ANCHORMONTH is outside range 1-12")
	}
	if ANCHORDAY < 1 || ANCHORDAY > 31 {
		return 0, 0, errors.New("value of ANCHORDAY is outside range 1-31")
	}

	anchor := time.Date(ANCHORYEAR, ANCHORMONTH, ANCHORDAY, 5, 0, 0, 0, time.Local)
	awYear, awWeek = anchor.ISOWeek()
	return awYear, awWeek, nil
}

// WeeksSinceAnchor returns how many weeks have elapsed between the anchor year and week
// and the current year and week
//
// It returns 0, error if the weeks since the anchor are less than 0
func WeeksSinceAnchor(cYear, cWeek, aYear, aWeek int) (wsAnchor int, err error) {

	// Add remaining weeks in anchor year
	wsAnchor += WeeksInYear(aYear) - aWeek

	// Add weeks so far in current year
	wsAnchor += cWeek

	// Add weeks for whole years
	if cYear > aYear+1 {
		for y := aYear + 1; y < cYear; y++ {
			wsAnchor += WeeksInYear(y)
		}
	}
	if wsAnchor < 0 {
		return 0, errors.New("value of wsAnchor less than 0")
	}

	return wsAnchor, nil
}

// IsBinCollected returns whether a bin is collected as a bool based on the given
// collection frequency and the number of weeks since the anchor week
func IsBinCollected(bin binData, wsAnchor int) bool {
	return wsAnchor%bin.Freq == 0
}

// weeksInYear returns how many weeks are in the given year, which varies based on
// whether it's a leap year (53) or a non leap year (52)
func WeeksInYear(year int) (d int) {
	t := time.Date(year, 12, 31, 0, 0, 0, 0, time.Local)
	d = t.YearDay()

	if d > 365 {
		return WEEKSINLEAPYEAR
	} else {
		return WEEKSINYEAR
	}
}
