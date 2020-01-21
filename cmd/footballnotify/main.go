package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/thang14/footballnotify/fire"
	"github.com/thang14/footballnotify/store"
	"github.com/thang14/footballnotify/types"
)

var footballAPIKey string

func getDBPath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return "./data"
	}
	return dbPath
}

func main() {
	s := store.New(getDBPath())
	footballAPIKey = s.GetFootballAPIKey()

	if footballAPIKey == "" {
		footballAPIKey = os.Getenv("API_KEY")
	}

	go watchEventChanges()

	http.HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body := types.Config{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				fmt.Fprintf(w, "parse json err: %s", err)
				return
			}

			if body.FootballAPIKey != "" {
				if err := s.SetFootballAPIKey(body.FootballAPIKey); err != nil {
					fmt.Fprintf(w, "set football api key error: %s", err)
					return
				}
				footballAPIKey = body.FootballAPIKey
			}
		}

		fmt.Fprintf(w, "ok")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func watchEventChanges() {
	fireService := fire.NewService()
	events := types.Events{}
	for {
		newEvents, err := getEvents()
		if err != nil {
			log.Printf("err: %s", err)
			continue
		}

		if len(events) == len(newEvents) {
			msgs := events.GetNotificationMessages(newEvents)
			fireService.SendMsgs(msgs)
		}

		if len(newEvents) > 0 {
			events = newEvents
		}

		// push message
		time.Sleep(5 * time.Second)
	}
}

func getEvents() (types.Events, error) {
	startTime := time.Now().Format("2006-01-02")
	endpoint := fmt.Sprintf("http://apiv2.apifootball.com/?action=get_events&APIkey=%s&from=%s&to=%s", footballAPIKey, startTime, startTime)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	events := make(types.Events, 0)
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}
