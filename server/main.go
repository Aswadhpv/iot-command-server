package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

var mqttClient mqtt.Client

func main() {
	// MQTT client options
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetClientID("command-server").
		SetCleanSession(true)

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connect error: %v", token.Error())
	}
	log.Println("Connected to MQTT broker")

	// HTTP router
	r := mux.NewRouter()
	r.HandleFunc("/devices/{id}/command", sendCommandHandler).Methods("POST")

	// Start HTTP server
	addr := ":8080"
	log.Printf("HTTP API listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

// Command represents the JSON payload for commands
type Command struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func sendCommandHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]

	var cmd Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	topic := fmt.Sprintf("devices/%s/commands", deviceID)
	payload, _ := json.Marshal(cmd)

	token := mqttClient.Publish(topic, 0, false, payload)
	token.WaitTimeout(5 * time.Second)
	if token.Error() != nil {
		http.Error(w, "failed to publish command", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
