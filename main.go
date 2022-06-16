package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Voting_phases
// type DBUSER struct {
// 	DiscordUserID string      `json:"userID"`
// 	VotingPhases  VotingPhase `json:"voting_phases"`
// }

type UserPoints struct {
	Uuid           string `json:"userID"`
	Points_phase_1 int    `json:"balance"`
}

type VotingPhase map[string]int8

type Leaderboard struct {
	Status string       `json:"status"`
	Points []UserPoints `json:"leaderboard"`
}

var (
	Url    = flag.String("url", "", "Points api url")
	Cookie = flag.String("cookie", "", "cookie of points api")
	Table  = flag.String("table", "", "Table name")
)

func init() {
	flag.Parse()
	prepareHttpClient()
	prepareRequest()
}

func prepareHttpClient() http.Client {
	api := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	return api
}

func prepareRequest() *http.Request {
	request, err := http.NewRequest(http.MethodGet, *Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Cookie", *Cookie)
	return request
}

func getJSONData(res *http.Response) []byte {
	jsonData, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return jsonData
}

func main() {

	api := prepareHttpClient()
	request := prepareRequest()

	res, getErr := api.Do(request)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	jsonData := getJSONData(res)

	var leaderboard Leaderboard
	jsonErr := json.Unmarshal([]byte(jsonData), &leaderboard)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//WRITE CSV
	file, err := os.Create("points.csv")

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, user := range leaderboard.Points {
		row := []string{user.Uuid, strconv.Itoa(user.Points_phase_1)}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
