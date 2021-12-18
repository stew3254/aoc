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

func getInput(scanner *bufio.Scanner) (output []string) {
	// Scan the input text
	for scanner.Scan() {
		// Fail on error
		err := scanner.Err()
		if err != nil {
			log.Fatalln(err)
		}

		// Use the text in the output
		line := scanner.Text()
		output = append(output, line)
	}
	return output
}

func problem1(input []string) (output int) {
	gamma := 0
	epsilon := 0
	length := len(input[0])
	for i := 0; i < length; i++ {
		zero := 0
		one := 0
		for _, diagnostic := range input {
			if diagnostic[i] == '0' {
				zero += 1
			} else {
				one += 1
			}
		}
		if one > zero {
			gamma += 1 << (length - i - 1)
		} else {
			epsilon += 1 << (length - i - 1)
		}
	}
	return gamma * epsilon
}

func filter(msb bool, pos int, input *[]string) {
	// Since we only have 1 string anyway just return
	if len(*input) == 1 {
		return
	}

	// Get the most common first bit
	zero := 0
	one := 0
	for _, diagnostic := range *input {
		if diagnostic[pos] == '0' {
			zero += 1
		} else {
			one += 1
		}
	}

	i := 0
	for {
		// Break when we go over length
		if i >= len(*input) {
			break
		}
		// Remove string if it doesn't belong
		if one >= zero && ((msb && (*input)[i][pos] != '1') || (!msb && (*input)[i][pos] != '0')) {
			*input = append((*input)[:i], (*input)[i+1:]...)
		} else if zero > one && ((msb && (*input)[i][pos] != '0') || (!msb && (*input)[i][pos] != '1')) {
			*input = append((*input)[:i], (*input)[i+1:]...)
		} else {
			// Bump up count if not removed
			i++
		}
	}
}

func problem2(input []string) (output int) {
	// Create the generator and scrubber lists
	generator := make([]string, len(input))
	scrubber := make([]string, len(input))

	// Copy them over
	copy(generator, input)
	copy(scrubber, input)

	length := len(input[0])

	for i := 0; i < length; i++ {
		filter(true, i, &generator)
		filter(false, i, &scrubber)
	}

	g, _ := strconv.ParseInt(generator[0], 2, 32)
	s, _ := strconv.ParseInt(scrubber[0], 2, 32)

	return int(g) * int(s)
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
	sol1 := problem1(input)
	log.Println("Solution 1:", sol1)
	sol2 := problem2(input)
	log.Println("Solution 2:", sol2)

	// Send the output to the server
	// client := makeClient()
	// submit(client, sol1)
	// submit(client, sol2)
}
