package main

import (
	"bufio"
	"container/list"
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

type Pos struct {
	X     int
	Y     int
	Value int
}

type Basins [][]Pos

func (b Basins) Len() int {
	return len(b)
}

func (b Basins) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Basins) Less(i, j int) bool {
	return len(b[i]) > len(b[j])
}

type Map [][]int

func (m Map) Width() int {
	if m.Height() == 0 {
		return 0
	}
	return len((m)[0])
}

func (m Map) Height() int {
	return len(m)
}

// Adjacent returns adjacent positions
func (m Map) Adjacent(x, y int) (adjacent []Pos) {
	if y-1 >= 0 {
		adjacent = append(adjacent, Pos{X: x, Y: y - 1, Value: m[y-1][x]})
	}
	if y+1 < m.Height() {
		adjacent = append(adjacent, Pos{X: x, Y: y + 1, Value: m[y+1][x]})
	}
	if x-1 >= 0 {
		adjacent = append(adjacent, Pos{X: x - 1, Y: y, Value: m[y][x-1]})
	}
	if x+1 < m.Width() {
		adjacent = append(adjacent, Pos{X: x + 1, Y: y, Value: m[y][x+1]})
	}
	return adjacent
}

func (m Map) Lows() (lows []Pos) {
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			// Go through adjacent positions
			isLow := true
			for _, pos := range m.Adjacent(x, y) {
				// If it's not smaller it's not a low point
				if m[y][x] >= pos.Value {
					isLow = false
					break
				}
			}
			// Add it if it's a low
			if isLow {
				lows = append(lows, Pos{X: x, Y: y, Value: m[y][x]})
			}
		}
	}
	return lows
}

func (m Map) Basins() (basins Basins) {
	lows := m.Lows()
	basins = make(Basins, 0, len(lows))
	// For each low find the basins
	for _, low := range lows {
		// Stack to hold surrounding
		surrounding := list.New()
		// Add the low to start with
		surrounding.PushFront(low)
		// Positions in the basin
		basin := make([]Pos, 0)
		seen := make(map[Pos]bool)
		// While the surrounding stack is not empty
		for surrounding.Len() > 0 {
			// Get the front of the list
			elem := surrounding.Front()
			pos := elem.Value.(Pos)
			for _, p := range m.Adjacent(pos.X, pos.Y) {
				// If we aren't at a peak and haven't seen this before, push onto the stack
				_, exists := seen[p]
				if !exists && p.Value < 9 {
					// Add to the stack
					surrounding.PushBack(p)
					// Mark as seen
					seen[p] = true
					// Add to the basin
					basin = append(basin, p)
				}
			}
			// Remove the position from the list
			surrounding.Remove(elem)
		}
		// Add the basin to the list of basins
		basins = append(basins, basin)
	}
	return basins
}

func getInput(scanner *bufio.Scanner) (output Map) {
	// Scan the input text
	for scanner.Scan() {
		// Fail on error
		err := scanner.Err()
		if err != nil {
			log.Fatalln(err)
		}

		// Use the text in the output
		line := scanner.Text()
		positions := make([]int, len(line))
		for i, c := range line {
			n, _ := strconv.Atoi(string(c))
			positions[i] = n
		}
		output = append(output, positions)
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

func problem1(input Map) (output int) {
	for _, pos := range input.Lows() {
		output += pos.Value + 1
	}
	return output
}

func problem2(input Map) (output int) {
	// Get the basins and sort from highest to lowest
	basins := input.Basins()
	sort.Sort(basins)

	output = 1
	count := 0
	for _, basin := range basins {
		if count == 3 {
			break
		}
		output *= len(basin)
		count++
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
