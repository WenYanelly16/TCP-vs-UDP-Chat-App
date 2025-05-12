# TCP vs UDP Chat App in Go

## 📚 Project Overview

This project explores the implementation and comparison of TCP and UDP protocols by building two versions of a concurrent chat application in Go. The goal is to examine the behavior, performance, and reliability of both protocols under various networking conditions.

---

## 🧩 Project Structure

```
TCP-vs-UDP/
├── pkg/
│   ├── message.go         # Common message structure and serialization
│   └── metrics.go         # Tracks latency, throughput, packet loss, etc.
│
├── TCP/
│   ├── client/
│   │   ├── client.go      # TCP client logic
│   │   └── main.go        # TCP client entry point
│   └── server/
│       ├── server.go      # TCP server logic
│       └── main.go        # TCP server entry point
│
├── UDP/
│   ├── client/
│   │   ├── client.go      # UDP client logic
│   │   └── main.go        # UDP client entry point
│   └── server/
│       ├── server.go      # UDP server logic
│       └── main.go        # UDP server entry point
│
├── test/
│   ├── integration_test.go  # End-to-end tests for TCP and UDP
│   └── performance_run.go   # Load simulation and metrics collection
│
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

---

## 🎯 Project Goals

- Build two chat systems using Go's `net` package:
  - One using **TCP**
  - One using **UDP**
- Clients can connect and exchange messages with a server.
- Messages are broadcasted to all connected clients.
- Handle client disconnection gracefully.
- Implement concurrency using goroutines and synchronization primitives.

---

## ⚙️ How to Run

### 🖥️ TCP Version

1. **Run TCP Server**
   go run TCP/server/main.go
2. Run **TCP Client(s)**
go run TCP/client/main.go

🌐 UDP Version
1. Run **UDP Server**
go run UDP/server/main.go
2. Run **UDP Client(s)**
go run UDP/client/main.go

📊 Metrics & Testing
**Collected Metrics**
- Latency
- Throughput
- Packet Loss

**Testing Features**
- Sudden client disconnection
- High message volume simulation
- Integration and performance testing under load (test/ folder)

To simulate network conditions and collect metrics:
go run test/performance_run.go

## 📽️ Deliverables
✅ Fully working TCP and UDP chat apps
✅ This README.md file
✅ PowerPoint presentation
✅ YouTube video demo

 ## 📈Performance Evaluation (To Be Presented)
Graphs and tables comparing:
- Latency under load
- Throughput with multiple clients
- Packet loss simulation

## Observations on:
- Protocol stability
- Real-time performance
- Scalability differences

 ## 🛠️Tools & Concepts Used
- Go’s net package for sockets
- Goroutines and channels for concurrency
- Mutexes and sync primitives
- Custom message struct and encoding
- Logging and metrics collection
- Structured project organization

## 👨‍🏫 Authors and Course Info
Course: CMPS2242 - Systems Programming & Computer Organization

Group Members: Wendy Alfaro & Irish Bigonia

Instructor: Mr. Dalwin Lewis

## Youtube Video Link:
https://youtu.be/bacj0fQx1qM

## Presentation Link:
https://www.canva.com/design/DAGmQ_0eyMs/oxbh1ORbbv0jSi3x4NNTnA/edit
