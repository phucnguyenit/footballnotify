package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thang14/footballnotify/fire"
	"github.com/thang14/footballnotify/types"
)

func main() {

	fire := fire.New()
	events := types.Events{}
	for {
		newEvents, err := getEvents()
		if err != nil {
			log.Printf("err: %s", err)
			continue
		}
		msgs := events.GetNotificationMessages(newEvents)
		fire.SendMsgs(msgs)
		events = newEvents

		// push message
		time.Sleep(5 * time.Second)
	}

}

func getEvents() (types.Events, error) {
	t := time.Now().Format("2006-01-02")
	endpoint := fmt.Sprintf("http://apiv2.apifootball.com/?action=get_events&APIkey=40f0117efa910900097069e35d39a7453fbbdb35af88a9072702e0690a15733d&from=%s&to=%s", t, t)
	log.Printf("get events: %s", endpoint)
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
