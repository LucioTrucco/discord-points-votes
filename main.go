package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Employee struct {
	ID  string
	Age int
}

type UserPoints struct {
	Uuid           string `json:"userID"`
	Points_phase_1 int    `json:"balance"`
}

type Leaderboard struct {
	Status      string       `json:"status"`
	Leaderboard []UserPoints `json:"leaderboard"`
}

var (
	Url    = flag.String("url", "", "Points api url")
	Cookie = flag.String("cookie", "", "cookie of points api")
)

func init() {
	flag.Parse()
}

func main() {

	url := Url

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", *Cookie)

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	jsonData, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var leaderboard Leaderboard
	jsonErr := json.Unmarshal([]byte(jsonData), &leaderboard)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	file, err := os.Create("points.csv")
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, user := range leaderboard.Leaderboard {
		row := []string{user.Uuid, strconv.Itoa(user.Points_phase_1)}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

}
