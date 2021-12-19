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

func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func abs(a int) int {
	if a > -1 {
		return a
	} else {
		return -a
	}
}

type Point struct {
	X int
	Y int
}

type Line struct {
	P1 Point
	P2 Point
}

type Grid map[Point]int

// Add the line to the grid
func (grid *Grid) Add(l Line, diagonal bool) {
	// Get point offsets
	x := l.P2.X - l.P1.X
	y := l.P2.Y - l.P1.Y

	// Get signs for direction
	xSign := 1
	ySign := 1

	if x < 0 {
		xSign = -1
	}
	if y < 0 {
		ySign = -1
	}

	for i := 0; i <= max(abs(x), abs(y)); i++ {
		var point Point
		if x != 0 && y != 0 {
			// Assume the diagonals have a slope of 1
			if diagonal {
				point = Point{X: l.P1.X + i*xSign, Y: l.P1.Y + i*ySign}
			} else {
				// Skip the rest of this and don't add the point
				continue
			}
		} else if x != 0 {
			point = Point{X: l.P1.X + i*xSign, Y: l.P1.Y}
		} else {
			point = Point{X: l.P1.X, Y: l.P1.Y + i*ySign}
		}
		// If the point exists on the grid, bump up intersection number. Otherwise, set to 1
		if _, exists := (*grid)[point]; exists {
			(*grid)[point] += 1
		} else {
			(*grid)[point] = 1
		}
	}
}

func getInput(scanner *bufio.Scanner) (output []Line) {
	// Scan the input text
	for scanner.Scan() {
		// Fail on error
		err := scanner.Err()
		if err != nil {
			log.Fatalln(err)
		}

		// Use the text in the output
		line := scanner.Text()
		pointsStr := strings.Split(line, " -> ")
		points := make([]Point, 0, 2)
		for _, point := range pointsStr {
			p := strings.Split(point, ",")
			x, _ := strconv.Atoi(p[0])
			y, _ := strconv.Atoi(p[1])
			points = append(points, Point{X: x, Y: y})
		}
		l := Line{P1: points[0], P2: points[1]}
		output = append(output, l)
	}
	return output
}

func problem1(input []Line) (output int) {
	grid := make(Grid)

	// Add all lines to the grid
	for _, line := range input {
		grid.Add(line, false)
	}

	// Get all points with intersections
	for _, v := range grid {
		if v > 1 {
			output++
		}
	}
	return output
}

func problem2(input []Line) (output int) {
	grid := make(Grid)

	// Add all lines to the grid
	for _, line := range input {
		grid.Add(line, true)
	}

	// Get all points with intersections
	for _, i := range grid {
		if i > 1 {
			output++
		}
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
