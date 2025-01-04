# Example [![GoDoc](https://godoc.org/github.com/thylong/example?status.png)](https://godoc.org/github.com/thylong/example) [![License](https://img.shields.io/badge/License-MIT%202.0-green.svg)](https://github.com/thylong/gonew-templates/blob/main/06-grpc-sqlc/LICENSE)
<!-- Logo -->
<!-- Labels (godoc, goreport, gocover, gosec, tests, doc link, Slack, license) -->
<!-- Pronunciation -->

<!-- Short description -->
Typical development workflow consist in the following steps:

1. Define your API contract with protocol buffers
2. Define your data model in a new file in `pkg/db/migrations` and write plain SQL to generate requests in `/pkg/db/queries`
3. Generate Go code for your proto & SQL binding with `make gen`
4. Write your server implementation in server Go package (including pgtype adaptations)

## Requirements

```bash
# For MacOS
brew install golang-migrate sqlc bufbuild/buf/buf grpcurl
```

## Quickstart

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/thylong/gonew-templates/06-grpc-sqlc github.com/example/myapp

```

### gRPCurl Example Queries for EventService

This document provides example `grpcurl` queries to interact with the `EventService` defined in the provided Protobuf schema. These examples assume your gRPC server is running locally on port `50051`.

---

#### **1. List All Services**

Use `grpcurl` to list all services available on the server:

```bash
grpcurl -plaintext localhost:50051 list
```

**Expected Output:**

```plaintext
events.EventService
```

---

#### **2. Describe a Service**

To describe the `EventService`, including its RPC methods:

```bash
grpcurl -plaintext localhost:50051 describe events.EventService
```

**Expected Output:**

```plaintext
events.EventService is a service:
  rpc DeleteEvent ( .events.DeleteEventRequest ) returns ( .events.DeleteEventResponse );
  rpc GetEvent ( .events.GetEventRequest ) returns ( .events.GetEventResponse );
  rpc GetEvents ( .events.GetEventsRequest ) returns ( .events.GetEventsResponse );
  rpc PutEvent ( .events.PutEventRequest ) returns ( .events.PutEventResponse );
```

---

#### **3. `GetEvents` Query**

Request a list of events with pagination and search filters:

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "search": "example"
}' localhost:50051 events.EventService/GetEvents
```

**Expected Output:**

```json
{
  "events": [
    {
      "event_privacy": "EVENT_PRIVACY_INTERNAL",
      "event_id": "123",
      "name": "Company All-Hands",
      "type": "Meeting",
      "department": "HR",
      "regions": "US",
      "tags": "monthly",
      "start_at": {
        "year": 2025,
        "month": 1,
        "day": 15,
        "hours": 10,
        "minutes": 30,
        "seconds": 0
      }
    }
  ],
  "total_count": 1,
  "page": 1,
  "page_size": 10
}
```

---

#### **4. `GetEvent` Query**

Request details for a specific event by `event_id`:

```bash
grpcurl -plaintext -d '{
  "event_id": "123"
}' localhost:50051 events.EventService/GetEvent
```

**Expected Output:**

```json
{
  "event": {
    "event_privacy": "EVENT_PRIVACY_INTERNAL",
    "event_id": "123",
    "name": "Company All-Hands",
    "type": "Meeting",
    "department": "HR",
    "regions": "US",
    "tags": "monthly",
    "start_at": {
      "year": 2025,
      "month": 1,
      "day": 15,
      "hours": 10,
      "minutes": 30,
      "seconds": 0
    }
  }
}
```

---

#### **5. `PutEvent` Query**

Create or update an event:

```bash
grpcurl -plaintext -d '{
  "event_privacy": "EVENT_PRIVACY_INTERNAL",
  "event_id": "123",
  "name": "Team Meeting",
  "type": "Meeting",
  "department": "Engineering",
  "regions": "EU",
  "tags": "sprint-planning",
  "start_at": {
    "year": 2025,
    "month": 1,
    "day": 18,
    "hours": 14,
    "minutes": 0,
    "seconds": 0
  }
}' localhost:50051 events.EventService/PutEvent
```

**Expected Output:**

```json
{
  "event": {
    "event_privacy": "EVENT_PRIVACY_INTERNAL",
    "event_id": "123",
    "name": "Team Meeting",
    "type": "Meeting",
    "department": "Engineering",
    "regions": "EU",
    "tags": "sprint-planning",
    "start_at": {
      "year": 2025,
      "month": 1,
      "day": 18,
      "hours": 14,
      "minutes": 0,
      "seconds": 0
    }
  }
}
```

---

#### **6. `DeleteEvent` Query**

Delete an event by `event_id`:

```bash
grpcurl -plaintext -d '{
  "eventID": "123"
}' localhost:50051 events.EventService/DeleteEvent
```

**Expected Output:**

```json
{}
```

#### **Query the health probe**

```bash

grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

## Benchmarks

## Installation

## Features

## FAQ

## License

## Contributing
