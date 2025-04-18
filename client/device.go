package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Command struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func main() {
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker URL")
	id := flag.String("id", "device123", "Device ID")
	flag.Parse()

	opts := mqtt.NewClientOptions().
		AddBroker(*broker).
		SetClientID(*id)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Connect error: %v", token.Error())
	}
	log.Printf("Device %s connected to %s", *id, *broker)

	topic := fmt.Sprintf("devices/%s/commands", *id)
	if token := client.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
		var cmd Command
		if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
			log.Printf("Bad command payload: %s", msg.Payload())
			return
		}
		log.Printf("Received command: %s %+v", cmd.Action, cmd.Params)
		// TODO: Insert real action handling here
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Subscribe error: %v", token.Error())
	}

	// Wait for Ctrl+C
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	client.Disconnect(250)
	log.Println("Device disconnected")
}
