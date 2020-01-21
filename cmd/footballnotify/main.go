package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/thang14/footballnotify/fire"
	"github.com/thang14/footballnotify/types"
)

func main() {
	watchEventChanges()
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
	apiKey := os.Getenv("API_KEY")
	endpoint := fmt.Sprintf("http://apiv2.apifootball.com/?action=get_events&APIkey=%s&from=%s&to=%s", apiKey, startTime, startTime)
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
