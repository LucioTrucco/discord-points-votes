package main

import (
	"encoding/csv"
	"encoding/json"
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

func main() {

	url := "https://points.city/api/guilds/961074073868308480/leaderboard"

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", "pointsID=xlTF88yqTNfUKr-NyD-nZdmImAk1S-7P.KCD3%2BTmVMC67bdOuYWa8SgyfhQ4P9c%2B7X8YtSYB43A0; Path=/; Expires=Sun, 17 Jul 2022 01:16:51 GMT; HttpOnly")

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
