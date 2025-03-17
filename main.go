package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/godbus/dbus/v5"
)

// NotificationHandler listens for notification actions
type NotificationHandler struct {
	conn *dbus.Conn
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() (*NotificationHandler, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %v", err)
	}

	return &NotificationHandler{
		conn: conn,
	}, nil
}

// SendNotification sends a notification with actions
func (h *NotificationHandler) SendNotification(appName, title, body string, actions []string) (uint32, error) {
	obj := h.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	// Convert actions to dbus format
	dbusActions := make([]string, len(actions)*2)
	for i, action := range actions {
		parts := strings.SplitN(action, "=", 2)
		if len(parts) == 2 {
			dbusActions[i*2] = parts[0]   // action key
			dbusActions[i*2+1] = parts[1] // action label
		}
	}

	fmt.Printf("%s", dbusActions)

	var id uint32
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		appName,                   // app_name
		uint32(0),                 // replaces_id
		"",                        // app_icon
		title,                     // summary
		body,                      // body
		dbusActions,               // actions
		map[string]dbus.Variant{}, // hints
		int32(-1),                 // timeout (-1 for default)
	)

	if err := call.Store(&id); err != nil {
		return 0, fmt.Errorf("failed to send notification: %v", err)
	}

	return id, nil
}

// Listen starts listening for notification action signals
func (h *NotificationHandler) Listen(actionHandlers map[string]func()) error {
	if err := h.conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.Notifications"),
		dbus.WithMatchMember("ActionInvoked"),
	); err != nil {
		return fmt.Errorf("failed to add match: %v", err)
	}

	c := make(chan *dbus.Signal, 10)
	h.conn.Signal(c)

actionLoop:
	for v := range c {
		log.Println(v.Name)
		if v.Name == "org.freedesktop.Notifications.ActionInvoked" {
			notificationID := v.Body[0].(uint32)
			actionKey := v.Body[1].(string)

			fmt.Printf("Notification %d: Action %s invoked\n", notificationID, actionKey)

			if handler, ok := actionHandlers[actionKey]; ok {
				handler()
				break actionLoop
			}
		}
	}

	return nil
}

func LogoutViaGnome() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %v", err)
	}

	obj := conn.Object("org.gnome.SessionManager", "/org/gnome/SessionManager")
	call := obj.Call("org.gnome.SessionManager.Logout", 0, uint32(0))

	return call.Err
}

// LogoutViaSystemd triggers a logout using loginctl
func LogoutViaSystemd() error {
	cmd := exec.Command("loginctl", "terminate-user", "$USER")
	return cmd.Run()
}

func HandleLogout() {
	fmt.Println("Logout action triggered!")

	// Try GNOME method first
	err := LogoutViaGnome()
	if err != nil {
		fmt.Printf("GNOME logout failed: %v\n", err)

		// Fall back to systemd method
		err = LogoutViaSystemd()
		if err != nil {
			fmt.Printf("Systemd logout failed: %v\n", err)
		}
	}
}

type Action struct {
	Title   string
	Trigger func()
}

type Actions []Action

func (ac Actions) Results() (map[string]func(), []string) {
	handlers := make(map[string]func())
	var actionString []string
	for _, action := range ac {
		functionName := GetFunctionName(action.Trigger)
		handlers[functionName] = action.Trigger
		actionString = append(actionString, fmt.Sprintf("%s=%s", functionName, action.Title))
	}
	return handlers, actionString
}

func HandleCall() {
	fmt.Println("Call action triggered!")
	// Implement your call functionality here
	// For example, you could launch a phone app or script
}

// HandleFind implements the find action
func HandleFind() {
	fmt.Println("Find action triggered!")
	// Implement your find functionality here
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func main() {
	handler, err := NewNotificationHandler()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	// Define action handlers
	actions := Actions{
		{
			Title:   "Logout Now",
			Trigger: HandleLogout,
		},
		{
			Trigger: HandleCall,
			Title:   "Call Me",
		},
	}

	actionHandlers, actionString := actions.Results()

	_, err = handler.SendNotification("hey", "title", "body", actionString)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
		if err != nil {
			return
		}
	}

	// Start listening for actions
	fmt.Println("Listening for notification actions...")
	if err := handler.Listen(actionHandlers); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error listening: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
