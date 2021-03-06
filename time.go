package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

// Allow scanning into date from using (*sql.Rows).Scan
func (t *Time) Scan(v interface{}) error {
	vt, err := time.Parse(time.RFC3339, string(v.(string)))
	if err != nil {
		return err
	}
	*t = Time(vt)
	return nil
}

func TimeFromShorthand(v string) (t time.Time) {
	layouts := []string{
		// Month
		"1",
		"01",
		// Year
		"2006",
		// Date for current year
		"01-02",
		// Date with year
		"20060102",
		"2006-01-02",
	}
	if v[len(v)-1] == 'd' {
		days, _ := strconv.Atoi(strings.TrimSuffix(v, "d"))
		t = time.Now().AddDate(0, 0, days)
		t = t.Round(time.Hour * 24)
		return t
	}
	if v[len(v)-1] == 'w' {
		weeks, _ := strconv.Atoi(strings.TrimSuffix(v, "w"))
		t = time.Now().AddDate(0, 0, weeks*7)
		t = t.Round(time.Hour * 24)
		return t
	}
	if v[len(v)-1] == 'm' {
		months, _ := strconv.Atoi(strings.TrimSuffix(v, "m"))
		t = time.Now().AddDate(0, months, 0)
		t = t.Round(time.Hour * 24)
		return t
	}
	if v[len(v)-1] == 'y' {
		years, _ := strconv.Atoi(strings.TrimSuffix(v, "y"))
		t = time.Now().AddDate(years, 0, 0)
		t = t.Round(time.Hour * 24)
		return t
	}
	for _, l := range layouts {
		if len(l) == len(v) {
			t, err := time.Parse(l, v)
			if err != nil {
				log.Fatal(err)
			}
			if t.Year() == 0 {
				t = t.AddDate(time.Now().Year(), 0, 0)
			}
			return t
		}
	}
	log.Fatalf("bad format provided: %s", v)
	return time.Time{}
}

func GetHours(d time.Duration) string {
	hours := d.Hours()
	s := ""
	if hours != 1 {
		s = "s"
	}
	return fmt.Sprintf("%.2f hour%s", hours, s)
}
