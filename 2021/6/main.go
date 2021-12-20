package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

var URL = "https://adventofcode.com"

func getAuth() *http.Cookie {
	data, err := ioutil.ReadFile("../../auth.txt")
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

func submit(client *http.Client, answer string) {
	// Get the year and day from the path
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	p, day := path.Split(dir)
	_, year := path.Split(p)

	// Submit the answers
	resp, err := client.PostForm(fmt.Sprintf("%s/%s/day/%s/answer", URL, year, day), url.Values{
		"level":  {day},
		"answer": {answer},
	})

	defer resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// See if it's the right answer or not
	if strings.Contains(doc.Find("main p").Text(), "not the right") {
		log.Println("Wrong answer")
	} else {
		log.Println("Correct!")
	}
}

type Cycle map[int]int

func getInput(scanner *bufio.Scanner) (output Cycle) {
	// Scan the input text
	for scanner.Scan() {
		// Fail on error
		err := scanner.Err()
		if err != nil {
			log.Fatalln(err)
		}

		// Use the text in the output
		line := scanner.Text()
		parts := strings.Split(line, ",")

		// Initialize the fish cycle counter
		output = make(Cycle)
		for _, p := range parts {
			f, _ := strconv.Atoi(p)
			output[f] += 1
		}
	}
	return output
}

func ageFish(cycle Cycle, days int) Cycle {
	last := 0
	// Age for the right number of days
	for i := 0; i < days; i++ {
		// Loop down through counter
		for day := 8; day >= 0; day-- {
			// Keep the current day, and set current day to last day and then replace last with current
			temp := cycle[day]
			cycle[day] = last
			last = temp
			if day == 0 {
				// Set day 6 to all of these fish
				cycle[6] += temp
				// Add a new lantern fish at day 8 as well
				cycle[8] = temp
			}
		}
	}
	return cycle
}

func totalFish(cycle Cycle) (total int) {
	for _, v := range cycle {
		total += v
	}
	return total
}

func problem1(input Cycle) (output int) {
	cycle := ageFish(input, 80)
	return totalFish(cycle)
}

func problem2(input Cycle) (output int) {
	cycle := ageFish(input, 256)
	return totalFish(cycle)
}

func main() {
	// Read in the input for the day
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	// Get the output to submit to the server
	input := getInput(scanner)

	// Copy the input for the next one
	input2 := make(Cycle)
	for k, v := range input {
		input2[k] = v
	}

	sol1 := problem1(input)
	log.Println("Solution 1:", sol1)
	sol2 := problem2(input2)
	log.Println("Solution 2:", sol2)

	// Send the output to the server
	// client := makeClient()
	// submit(client, sol1)
	// submit(client, sol2)
}
