package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// alertGroup is the data read from a webhook call
type alertGroup struct {
	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`

	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   alerts `json:"alerts"`

	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// alerts is a slice of Alert
type alerts []alert

// alert holds one alert for notification templates.
type alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
}

// parse gets a webhook payload and parses it returning a prometheus
// template.Data object if successful
func parse(payload []byte) (*alertGroup, error) {
	d := alertGroup{}
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json webhook payload: %w", err)
	}
	return &d, nil
}

// printJson prints all alerts in json format
func printJson(ag *alertGroup, m *sync.Mutex) {
	m.Lock()
	for _, alert := range ag.Alerts {
		out := map[string]string{"status": alert.Status}

		for k, v := range alert.Labels {
			out[k] = v
		}
		for k, v := range alert.Annotations {
			out[k] = v
		}
		out["startsAt"] = alert.StartsAt.Truncate(time.Millisecond).String()
		out["endsAt"] = alert.EndsAt.Truncate(time.Millisecond).String()

		jout, err := json.Marshal(out)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", jout)
	}
	m.Unlock()
}

// printKV iterates over the alertgroup and prints all alerts as key value pairs
func printKV(ag *alertGroup, m *sync.Mutex) {
	m.Lock()
	for _, alert := range ag.Alerts {

		fmt.Printf("\"status: %s\", ", alert.Status)

		for k, v := range alert.Labels {
			fmt.Printf("\"%s: %s\", ", k, v)
		}
		for k, v := range alert.Annotations {
			fmt.Printf("\"%s: %s\", ", k, v)
		}
		fmt.Printf("\"startsAt: %s\", \"endsAt: %s\"\n", alert.StartsAt.Truncate(time.Millisecond), alert.EndsAt.Truncate(time.Millisecond))
	}
	m.Unlock()
}

func main() {
	var m sync.Mutex

	jsonOutput := false
	if os.Getenv("JSON_OUTPUT") == "true" {
		jsonOutput = true
	}

	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		ag, err := parse(b)
		if err != nil {
			panic(err)
		}

		if jsonOutput {
			printJson(ag, &m)
		} else {
			printKV(ag, &m)
		}
	})

	server := &http.Server{
		Addr:              ":5001",
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}
