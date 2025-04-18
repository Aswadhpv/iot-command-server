# IoT Command Server

This project implements a simple IoT command-and-control server in Go, using MQTT for device communication and a RESTful HTTP API to send commands to devices.

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Project Structure](#project-structure)
4. [Installation & Setup](#installation--setup)
   - [Clone the Repository](#clone-the-repository)
   - [Install Dependencies](#install-dependencies)
   - [Install MQTT Broker](#install-mqtt-broker)
5. [Configuration](#configuration)
6. [Usage](#usage)
   - [Run the Server](#run-the-server)
   - [Run a Device Client](#run-a-device-client)
   - [Send Commands](#send-commands)
7. [Testing](#testing)
8. [Enhancements & Next Steps](#enhancements--next-steps)
9. [Contributing](#contributing)
10. [License](#license)

---

## Overview

This Go-based IoT Command Server allows you to send JSON-formatted commands from an HTTP API endpoint to one or more IoT clients over MQTT. Clients subscribe to specific topics and execute or log the commands they receive.

## Prerequisites

- **Go** (version 1.18 or newer)
- **MQTT Broker** (e.g., Mosquitto)
- **Git** (to clone the repository)
- **Docker** (optional, if you prefer running Mosquitto in a container)

## Project Structure

```
iot-command-server/
├── server/
│   └── main.go       # HTTP API + MQTT publisher
├── client/
│   └── device.go     # Sample IoT client subscribing to commands
├── go.mod
├── go.sum
└── README.md         # Project documentation
```

## Installation & Setup

### Clone the Repository

```bash
git clone https://github.com/Aswadhpv/iot-command-server.git
cd iot-command-server
```

### Install Dependencies

```bash
go mod download
```

### Install MQTT Broker

#### Windows

- **Chocolatey**:  `choco install mosquitto -y`
- **Scoop**:      `scoop install mosquitto`
- **Manual MSI**:
  1. Download the MSI from https://mosquitto.org/download/
  2. Run the installer and add `C:\Program Files\mosquitto` to your PATH.
- **Docker**:     `docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto`

#### Linux/macOS

```bash
# Debian/Ubuntu
tsudo apt update && sudo apt install mosquitto -y
# macOS (Homebrew)
brew install mosquitto
```

## Configuration

By default, the server and client connect to `tcp://localhost:1883`. You can override this by setting the `MQTT_BROKER_URL` environment variable:

```bash
# Example
env MQTT_BROKER_URL="tcp://broker.example.com:1883"
```

## Usage

### Run the Server

```bash
cd server
go run main.go
```

The server will:

- Connect to the MQTT broker
- Expose `POST /devices/{id}/command`
- Publish any received JSON payload to topic `devices/{id}/commands`

### Run a Device Client

```bash
cd client
go run device.go --id=device123
```

The client will:

- Connect to the same MQTT broker
- Subscribe to `devices/device123/commands`
- Log any commands it receives

### Send Commands

Use `curl`, Postman, or any HTTP client:

```bash
curl -X POST http://localhost:8080/devices/device123/command \
  -H "Content-Type: application/json" \
  -d '{"action":"toggle","params":{"pin":13}}'
```

Expected client output:

```
Received command: toggle map[pin:13]
```

## Testing

- Start multiple clients with different IDs to verify concurrent handling.
- Try malformed JSON or missing fields to test error responses.
- Secure the endpoint with a simple API key header and test unauthorized requests.

## Enhancements & Next Steps

- **Security**: Add TLS for MQTT and HTTPS for the API, plus authentication (JWT/API keys).
- **Device Registry**: Persist registered clients in a database (e.g., PostgreSQL).
- **Telemetry**: Allow clients to publish sensor data back to `devices/{id}/telemetry`.
- **Web UI**: Build a dashboard to list devices and send commands via browser.
- **Docker Compose**: Orchestrate server, broker, and database in a single file.

