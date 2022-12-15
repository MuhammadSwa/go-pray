package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Response []struct {
	Data []struct {
		Timings map[string]string `json:"timings"`
		Date    struct {
			Hijri struct {
				Date string `json:"date"`
			} `json:"hijri"`
		} `json:"date"`
	} `json:"data"`
}

func (r *Response) fetch(i int, conf ConfFile) []byte {
	var (
		city     = conf.City
		country  = conf.Country
		method   = conf.Method
		thisYear = time.Now().Year()
		iS       = fmt.Sprintf("%02d", i)
		// api for the month
		// city(string) - country (string)- method (int) month (int) - year (int)
		API = fmt.Sprintf("http://api.aladhan.com/v1/calendarByCity?city=%s&country=%s&method=%d&month=%s&year=%d",
			city, country, method, iS, thisYear)
	)

	fmt.Printf("Fetching data for %s\n", time.Month(i))

	resp, err := http.Get(API)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

func (r *Response) save(dataPath string, i int, body []byte) {
	thisMonth := int(time.Now().Month())
	file, _ := os.OpenFile(dataPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	// construct a json array [month1,month2...]
	if i == thisMonth {
		fmt.Fprintln(file, "[")
	}

	file.Write(body)

	if i != 12 {
		fmt.Fprintln(file, ",")
	}

	if i == 12 {
		fmt.Fprintln(file, "]")
	}

}

func (r *Response) FetchAndSaveData(conf ConfFile) {
	var (
		cnt       = 0
		thisYear  = time.Now().Year()
		thisMonth = int(time.Now().Month())
		dataPath  = conf.DataPath
	)
	fmt.Printf("Fetching data for year %d...\n", thisYear)

	// loop from this month to the end of the year
	for i := thisMonth; i <= 12; i++ {
		body := r.fetch(i, conf)

		r.save(dataPath, i, body)

		cnt++
		// sleep for each request so the api doesn't return an error
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("Fetched data for %d months\n", cnt)
}

func (r *Response) ShouldIFetch(conf ConfFile) {
	// instead of this using openFile(os.o_Create)
	stat, err := os.Stat(conf.DataPath)
	if err != nil {
		r.FetchAndSaveData(conf)
		return
	}

	if stat.ModTime().Year() == time.Now().Year() {
		return
	}
	r.FetchAndSaveData(conf)
}

func (r *Response) UnMarshal(dataPath string) {
	file, err := os.ReadFile(dataPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	err = json.Unmarshal(file, r)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
