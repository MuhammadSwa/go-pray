package main

import (
	"fmt"
	"time"
)

// TODO add hijri date option

var prayersOrder = [...]string{"Fajr", "Sunrise", "Dhuhr", "Asr", "Maghrib", "Isha"}

func parseTimings(prayers map[string]string) map[string]time.Time {
	prayersParsed := make(map[string]time.Time, 6)

	for key, value := range prayers {
		parsedTime, _ := time.Parse("15:04 (EET)", value)
		prayersParsed[key] = parsedTime
	}
	return prayersParsed
}

// List all prayers in order
func ListTimings(prayers map[string]time.Time) {
	for _, name := range prayersOrder {
		prayerName := prayers[name].Format("15:04")
		fmt.Printf("%s %s\n", name, prayerName)
	}
}

// Get the time left for next prayer
func NextPrayer(p map[string]time.Time, order [6]string) (string, time.Duration) {
	var (
		nextPrayer string
		timeLeft   time.Duration
	)

	now := time.Now().Format("15:04:05")
	nowParsed, _ := time.Parse("15:04:05", now)

	for _, name := range order {
		if p[name].After(nowParsed) {
			nextPrayer = name
			timeLeft = p[name].Sub(nowParsed)
			break
		} else {
			// It's fajr next day
			nextPrayer = order[0]
			fajrNextDay := p["Fajr"].AddDate(0, 0, 1)
			timeLeft = fajrNextDay.Sub(nowParsed)
			continue
		}
	}

	return nextPrayer, timeLeft
}
