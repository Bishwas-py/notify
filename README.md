# Go Notification Library

A simple, lightweight library for sending desktop notifications with action buttons and sound alerts in Go. Uses D-Bus for Linux desktop environments.

## Installation

```
go get github.com/Bishwas-py/notify
```

## Features

- Send notifications with custom actions
- Add sound notifications using standard freedesktop sound names or custom sound files
- Handle action button clicks with custom triggers
- Set notification timeout
- Clean API with struct-based configuration

## Basic Usage

```go
package main

import (
	"github.com/Bishwas-py/notify"
	"log"
	"time"
)

func HandleClick() {
	log.Println("Button clicked!")
}

func main() {
	// Create notification with an action button
	notification := notify.Notification{
		AppID:   "example",
		Title:   "Hello World",
		Body:    "This is a notification",
		Timeout: int(10 * time.Second),
		Actions: notify.Actions{
			{
				Title:   "Click Me",
				Trigger: HandleClick,
			},
		},
	}

	// Display notification and wait for response
	_, err := notification.Trigger()
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
```

## Sound Notifications

The library supports both standard freedesktop sound names and custom sound files:

```go
// Using a standard sound
notification := notify.Notification{
	Title:   "Alarm",
	Body:    "Wake up!",
	Timeout: int(5 * time.Second),
}
notification.SetSoundByName(notify.AlarmClockElapsed)

// Or using a custom sound file
notification := notify.Notification{
	Title:   "Custom Sound",
	Body:    "This has a custom sound",
	Timeout: int(5 * time.Second),
}
notification.SetSoundByPath("/path/to/your/sound.wav")
```

## Available Sound Names

The library provides constants for standard freedesktop sound names:

- `AlarmClockElapsed`
- `Bell`
- `CameraShutter`
- `Complete`
- `DeviceAdded`
- `DeviceRemoved`
- `DialogError`
- `DialogInformation`
- `DialogWarning`
- `Message`
- `MessageNewInstant`
- `PhoneIncomingCall`
- `PowerPlug`
- `PowerUnplug`
- `ScreenCapture`
- `ServiceLogin`
- `ServiceLogout`
- `TrashEmpty`
- `WindowAttention`
- And more...

## Multiple Actions Example

```go
actions := notify.Actions{
	{
		Title:   "Accept",
		Trigger: HandleAccept,
	},
	{
		Title:   "Decline",
		Trigger: HandleDecline,
	},
}

notification := notify.Notification{
	AppID:   "example",
	Title:   "Meeting Invitation",
	Body:    "You're invited to a meeting at 3 PM",
	Timeout: int(30 * time.Second),
	Actions: actions,
}

// Add a sound alert
notification.SetSoundByName(notify.MessageNewInstant)

// Trigger the notification
notificationID, err := notification.Trigger()
```

## Implementation Details

This library comprises three main components:

- `NotificationHandler`: Core D-Bus communication layer for sending/receiving notifications
- `Action`: Defines clickable buttons with associated handler functions
- `Notification`: High-level struct that configures notification parameters

The implementation uses Go's channels and select statements to handle notification responses with timeouts. Function names are obtained through reflection to match D-Bus action signals with corresponding handlers. The global `notificationHandler` variable is initialized during package import through the `init()` function.

## Requirements

- Linux desktop environment with D-Bus support
- Go 1.13 or higher

## License

[MIT](LICENSE)