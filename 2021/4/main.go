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

type Tile struct {
	Value int
	Drawn bool
}

type Board struct {
	Size  int
	State [][]Tile
}

// Play sets all tiles on a board with that value
func (b *Board) Play(value int) {
	for row := 0; row < b.Size; row++ {
		for col := 0; col < b.Size; col++ {
			if b.State[row][col].Value == value {
				b.State[row][col].Drawn = true
			}
		}
	}
}

// Win checks if the board is in a winning state
func (b *Board) Win() bool {
	colWin := make([]int, b.Size)
	for i := 0; i < b.Size; i++ {
		rowWin := 0
		for j := 0; j < b.Size; j++ {
			if b.State[i][j].Drawn {
				rowWin += 1
				colWin[j] += 1
			}
		}
		if rowWin == b.Size {
			return true
		}
	}
	for _, col := range colWin {
		if col == b.Size {
			return true
		}
	}
	return false
}

// Score returns the board score
func (b *Board) Score(drawn int) int {
	total := 0
	for _, row := range b.State {
		for _, tile := range row {
			// Add tile if not drawn
			if !tile.Drawn {
				total += tile.Value
			}
		}
	}
	return total * drawn
}

func getInput(scanner *bufio.Scanner) (draws []int, output []*Board) {
	// Get the first line input
	scanner.Scan()
	drawsStr := strings.Split(scanner.Text(), ",")
	draws = make([]int, len(drawsStr))
	for i, s := range drawsStr {
		draw, _ := strconv.Atoi(s)
		draws[i] = draw
	}
	// Skip the next line
	scanner.Scan()

	// Scan the input text
	board := &Board{}
	count := 0
	for scanner.Scan() {
		// Fail on error
		err := scanner.Err()
		if err != nil {
			log.Fatalln(err)
		}

		// Use the text in the output
		line := scanner.Text()
		parts := strings.Fields(line)

		// See if we are on a blank line
		for len(parts) == 0 {
			// Add the old board to the output
			output = append(output, board)
			// Make a new board
			board = &Board{}
			// Reset counter
			count = 0
			// Skip the line
			if !scanner.Scan() {
				return draws, output
			}
			line = scanner.Text()
			parts = strings.Fields(line)
		}

		// Set board size and make a board
		if board.Size == 0 {
			board.Size = len(parts)
			board.State = make([][]Tile, board.Size)
		}

		tiles := make([]Tile, board.Size)
		for i, s := range parts {
			v, _ := strconv.Atoi(s)
			tiles[i].Value = v
			tiles[i].Drawn = false
		}
		// Set the board row output
		board.State[count] = tiles
		count++
	}
	output = append(output, board)
	return draws, output
}

func problem1(draws []int, boards []*Board) (output int) {
	// Keep drawing until a win
	for _, draw := range draws {
		for _, board := range boards {
			// Play the tile
			board.Play(draw)
			// If this board wins, return its score
			if board.Win() {
				return board.Score(draw)
			}
		}
	}

	return output
}

func problem2(draws []int, boards []*Board) (output int) {
	winners := make(map[*Board]bool)
	// Keep drawing until a win
	for _, draw := range draws {
		for _, board := range boards {
			// If this board wins, return its score
			if _, exists := winners[board]; !exists {
				// Play the tile
				board.Play(draw)

				if board.Win() {
					winners[board] = true
				}
			}

			// Print the score of the last winning board
			if len(winners) == len(boards) {
				return board.Score(draw)
			}
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
	draws, boards := getInput(scanner)
	sol1 := problem1(draws, boards)
	log.Println("Solution 1:", sol1)
	sol2 := problem2(draws, boards)
	log.Println("Solution 2:", sol2)

	// Send the output to the server
	// client := makeClient()
	// submit(client, sol1)
	// submit(client, sol2)
}
