package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Logging struct {
	Timestamp string  `json:"timestamp"`
	Device    string  `json:"device_id"`
	CO        float64 `json:"co"`
	Humidity  float64 `json:"humid"`
	Light     bool    `json:"light"`
	LPG       float64 `json:"lpg"`
	Motion    bool    `json:"motion"`
	Smoke     float64 `json:"smoke"`
	Temp      float64 `json:"temp"`
}

func main() {
	parsedCsv := parseCsv("simulation_data/" + os.Getenv("DEVICE") + ".csv")
	endpoint := os.Getenv("ENDPOINT")

	for {
		for _, data := range parsedCsv {
			sendRequest(endpoint, data)
			time.Sleep(1 * time.Second)
		}
	}
}

func parseCsv(filename string) []Logging {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("OpenFile fatal error: %v", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	iotData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("ReadAll fatal error: %v", err)
	}

	var parsedData []Logging
	for i, iot := range iotData {
		if i == 0 {
			continue
		}

		var currData Logging
		currData.Timestamp = iot[0]
		currData.Device = iot[1]
		currData.CO, _ = strconv.ParseFloat(iot[2], 64)
		currData.Humidity, _ = strconv.ParseFloat(iot[3], 64)
		currData.Light, _ = strconv.ParseBool(iot[4])
		currData.LPG, _ = strconv.ParseFloat(iot[5], 64)
		currData.Motion, _ = strconv.ParseBool(iot[4])
		currData.Smoke, _ = strconv.ParseFloat(iot[7], 64)
		currData.Temp, _ = strconv.ParseFloat(iot[8], 64)

		parsedData = append(parsedData, currData)
	}

	return parsedData
}

func sendRequest(endpoint string, data Logging) {
	body, _ := json.Marshal(data)
	log.Println(string(body))
	bodyBuffer := bytes.NewBuffer(body)

	res, err := http.Post(endpoint, "application/json", bodyBuffer)
	if err != nil {
		log.Fatalf("Request fatal error: %v", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	log.Println(string(resBody))
}
