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
	"sort"
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

func abs(a int) int {
	if a > -1 {
		return a
	} else {
		return -a
	}
}

func sum(args ...int) (total int) {
	for _, v := range args {
		total += v
	}
	return total
}

func product(args ...int) (total int) {
	for _, v := range args {
		total *= v
	}
	return total
}

func max(args ...int) (m int) {
	if len(args) == 0 {
		return 0
	}
	temp := args[0]
	for _, v := range args[1:] {
		if v > temp {
			temp = v
			m = temp
		}
	}
	return m
}

func min(args ...int) (m int) {
	if len(args) == 0 {
		return 0
	}
	temp := args[0]
	for _, v := range args[1:] {
		if v < temp {
			temp = v
			m = temp
		}
	}
	return m
}

func average(args ...int) float32 {
	return float32(sum(args...)) / float32(len(args))
}

func mostCommon(args ...int) (common []int) {
	m := make(map[int]int)
	// Count up the number of occurrences of each number
	for _, v := range args {
		m[v] += 1
	}

	value := 0
	for k, v := range m {
		// If the values are the same append the key
		if value == v {
			common = append(common, k)
		}

		temp := max(value, v)
		// If it was bigger, add just that key
		if temp > value {
			value = temp
			common = []int{k}
		}
	}

	return common
}

func median(args ...int) int {
	sort.Ints(args)

	middle := len(args) / 2

	// Check if odd
	if len(args)&1 == 1 {
		return args[middle]
	}

	return (args[middle-1] + args[middle]) / 2
}

func sequence(n int) (total int) {
	return n * (n + 1) / 2
}

func problem1(input []string) (output int) {
	return output
}

func problem2(input []string) (output int) {
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
