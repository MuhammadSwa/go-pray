package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	prayersParsed, HijriDate := run()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			ListTimings(prayersParsed)

		case "next":
			prayerName, timeLeft := NextPrayer(prayersParsed, prayersOrder)
			fmt.Printf("%v left for %s\n", timeLeft, prayerName)

		case "date":
			fmt.Println(HijriDate)

		case "help":
			help()
		default:
			fmt.Fprintln(os.Stderr, "Please provide a valid argument")
			help()
		}
	} else {
		help()
	}
}
func run() (map[string]time.Time, string) {
	// load the config file
	conf := ConfFile{}
	conf.GetConf()

	data := Response{}
	data.ShouldIFetch(conf)
	data.UnMarshal(conf.DataPath)

	// find index of month in and day in Json file
	m := int(time.Now().Month())
	l := len(data)
	indexOfMonth := (m - (12 - l)) - 1
	indexOfDay := time.Now().Day() - 1

	dataOfToday := data[indexOfMonth].Data[indexOfDay]
	dayTimings := dataOfToday.Timings
	HijriDate := dataOfToday.Date.Hijri.Date

	prayersParsed := parseTimings(dayTimings)

	return prayersParsed, HijriDate
}
