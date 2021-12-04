package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var URL = "https://adventofcode.com"

func getAuth() *http.Cookie {
	data, err := ioutil.ReadFile("auth.txt")
	if err != nil {
		log.Fatalln(err)
	}

	return &http.Cookie{
		Name:     "session",
		Value:    string(data),
		Path:     "/",
		Domain:   ".adventofcode.com",
		Secure:   true,
		HttpOnly: true,
	}
}

func makeClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}
	urlObj, _ := url.Parse(URL)
	jar.SetCookies(urlObj, []*http.Cookie{getAuth()})

	return &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       0,
	}
}

func getDocument(client *http.Client, url string) (*goquery.Document, error) {
	// Get the HTML
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, errors.New(resp.Status)
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}
func getYear() (year string) {
	if len(os.Args) < 2 {
		// Get current year
		year = strconv.Itoa(time.Now().Year())
	} else {
		// Assume this is a valid year
		year = os.Args[1]
	}
	return year
}

func getDays(client *http.Client, year string) []string {
	doc, err := getDocument(client, URL+"/"+year)
	if err != nil {
		log.Fatalln(err)
	}
	dayFields := doc.Find(".calendar a")
	dayNum := dayFields.Length()
	days := make([]string, 0, dayNum)
	dayFields.Each(func(i int, s *goquery.Selection) {
		// Add the day number to the slice
		days = append(days, strings.TrimSpace(s.Find(".calendar-day").Text()))
	})
	return days
}

func initializeDays(client *http.Client, year string, days []string) {
	for _, day := range days {
		inputName := fmt.Sprintf("%s/%s/input.txt", year, day)
		goName := fmt.Sprintf("%s/%s/main.go", year, day)
		// See if the file exists
		if _, err := os.Stat(inputName); err == nil {
			continue
		}

		// Get the input for the day
		urlStr := fmt.Sprintf("%s/%s/day/%s/input", URL, year, day)
		resp, err := client.Get(urlStr)
		if err != nil {
			log.Fatalln(err)
		}

		if resp.StatusCode != 200 {
			log.Println(urlStr, resp.Status)
			continue
		}

		outFile, err := os.Create(inputName)
		if err != nil {
			log.Fatalln(err)
		}

		io.Copy(outFile, resp.Body)
		outFile.Close()

		// Copy the template.go into the new directory
		copyFile(goName, "template.go")
	}
}

func copyFile(dst, src string) {
	// Read the source file
	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatalln(err)
	}
	defer srcFile.Close()

	// Create the file if it doesn't exist
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Fatalln(err)
	}
	defer dstFile.Close()

	// Copy the data over
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Sync the file contents to the disk
	err = dstFile.Sync()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	client := makeClient()

	// Get the year
	year := getYear()
	// Get the days of this year
	days := getDays(client, year)

	// Create the directory for each day if it doesn't exist
	for _, day := range days {
		err := os.MkdirAll(year+"/"+day, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Get all inputs if they don't exist
	initializeDays(client, year, days)
}
