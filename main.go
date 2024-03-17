package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func randomValue() int {
	return rand.Intn(100) + 1
}

func updateJSONFile() {
	for {
		status := Status{
			Water: randomValue(),
			Wind:  randomValue(),
		}

		data, _ := json.MarshalIndent(status, "", "  ")
		file, err := os.Create("status.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Printf("Status updated: Water %d, Wind %d\n", status.Water, status.Wind)

		time.Sleep(15 * time.Second)
	}
}

func main() {
	go updateJSONFile()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("status.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var status Status
		err = json.NewDecoder(file).Decode(&status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var statusWater, statusWind string
		if status.Water < 5 {
			statusWater = "Aman"
		} else if status.Water >= 6 && status.Water <= 8 {
			statusWater = "Siaga"
		} else {
			statusWater = "Bahaya"
		}

		if status.Wind < 6 {
			statusWind = "Aman"
		} else if status.Wind >= 7 && status.Wind <= 15 {
			statusWind = "Siaga"
		} else {
			statusWind = "Bahaya"
		}

		templateHTML := template.Must(template.ParseFiles("index.html"))

		dataMap := map[string]interface{}{
			"Water":       status.Water,
			"Wind":        status.Wind,
			"StatusWater": statusWater,
			"StatusWind":  statusWind,
		}

		err = templateHTML.Execute(w, dataMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
