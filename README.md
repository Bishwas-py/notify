# Go Notification Library

A simple, lightweight library for sending desktop notifications with action buttons in Go. Uses D-Bus for Linux desktop
environments.

## Installation

```
go get github.com/Bishwas-py/notify
```

## Features

- Send notifications with custom actions
- Handle action button clicks
- Set notification timeout
- Clean API with struct-based configuration

## Usage

```go
notification := Notification{
 Title:   "Hello World",
 Body:    "This is a notification",
 Timeout: int(10 * time.Second),
 Actions: Actions{{Title: "Click Me", Trigger: HandleClick}},
}

notification.Trigger() // Displays notification and waits for response
```

## Implementation Details

This library comprises three main components:

- `NotificationHandler`: Core D-Bus communication layer for sending/receiving notifications
- `Action`: Defines clickable buttons with associated handler functions
- `Notification`: High-level struct that configures notification parameters

The implementation uses Go's channels and select statements to handle notification responses with timeouts. Function
names are obtained through reflection to match D-Bus action signals with corresponding handlers. The global
`notificationHandler` variable is initialized during package import through the `init()` function.
