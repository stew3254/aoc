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

func advent(scanner *bufio.Scanner) (output string) {
	// Scan the input text
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
		// Use the text in the output
		line := scanner.Text()
		output += line + "\n"
	}

	return output
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
	output := advent(scanner)
	log.Println(output)

	// Send the output to the server
	// client := makeClient()
	// submit(client, output)
}
