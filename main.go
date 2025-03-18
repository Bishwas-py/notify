package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"time"
)

// NotificationHandler listens for notification actions
type NotificationHandler struct {
	conn *dbus.Conn
}

var notificationHandler *NotificationHandler

func init() {
	log.Printf("Running on Go %s", runtime.Version())

	_, err := exec.LookPath("loginctl")
	if err != nil {
		log.Printf("Warning: loginctl not found in PATH. Systemd logout method may not work.")
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Printf("Warning: D-Bus session bus not available: %v", err)
		log.Printf("Notifications may not work properly")
	}

	notificationHandler = &NotificationHandler{
		conn: conn,
	}
}

// SendNotification sends a notification with actions
func (h *NotificationHandler) SendNotification(appName, title, body string, dbusActions []string, timeout int) (uint32, error) {
	obj := h.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var id uint32
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		appName,                   // app_name
		uint32(0),                 // replaces_id
		"",                        // app_icon
		title,                     // summary
		body,                      // body
		dbusActions,               // actions
		map[string]dbus.Variant{}, // hints
		timeout,
	)

	if err := call.Store(&id); err != nil {
		return 0, fmt.Errorf("failed to send notification: %v", err)
	}

	return id, nil
}

// Listen starts listening for notification action signals
func (h *NotificationHandler) Listen(actionHandlers map[string]func(), notificationID uint32, timeout int) error {
	if err := h.conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.Notifications"),
		dbus.WithMatchMember("ActionInvoked"),
	); err != nil {
		return fmt.Errorf("failed to add match: %v", err)
	}

	c := make(chan *dbus.Signal, 10)
	h.conn.Signal(c)

	timer := time.NewTimer(time.Duration(timeout))
	defer timer.Stop()

actionLoop:
	for {
		log.Println("Waiting for signal...")
		select {
		case v, ok := <-c:
			if !ok {
				break actionLoop
			}

			if v.Name == "org.freedesktop.Notifications.ActionInvoked" {
				curNid := v.Body[0].(uint32) // Current Notification ID
				actionKey := v.Body[1].(string)

				log.Printf("Notification %d: Action %s invoked %d\n", curNid, actionKey, notificationID)
				if curNid != notificationID {
					continue
				}
				if handler, ok := actionHandlers[actionKey]; ok {
					handler()
					_ = h.conn.Close()
					break actionLoop
				}
			}
			if v.Name == "org.freedesktop.Notifications.NotificationClosed" {
				break actionLoop
			}

		case <-timer.C:
			println("Timeout reached, breaking loop")
			break actionLoop
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
	log.Println("Logout action triggered!")

	// Try GNOME method first
	err := LogoutViaGnome()
	if err != nil {
		log.Printf("GNOME logout failed: %v\n", err)

		// Fall back to systemd method
		err = LogoutViaSystemd()
		if err != nil {
			log.Printf("Systemd logout failed: %v\n", err)
		}
	}
}

type Action struct {
	Title   string
	Trigger func()
}

type Actions []Action

func (ac Actions) Results() (map[string]func(), []string) {
	handlers := make(map[string]func(), len(ac))
	actions := make([]string, 0, len(ac)*2)

	for _, action := range ac {
		functionName := GetFunctionName(action.Trigger)
		handlers[functionName] = action.Trigger
		actions = append(actions, functionName, action.Title)
	}
	return handlers, actions
}

func HandleCall() {
	log.Println("Call action triggered!")
	// Implement your call functionality here
	// For example, you could launch a phone app or script
}

// HandleFind implements the find action
func HandleFind() {
	log.Println("Find action triggered!")
	// Implement your find functionality here
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func main() {
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

	listenHandlers, dbusActions := actions.Results()
	timeout := 10 * time.Second

	notificationID, err := notificationHandler.SendNotification("hey", "title", "body", dbusActions, int(timeout))
	log.Printf(" ID: %d, Error: %v", notificationID, err)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return
	}

	log.Printf("Listening to notification id: %d", notificationID)
	if err := notificationHandler.Listen(listenHandlers, notificationID, int(timeout)); err != nil {
		log.Printf("Failed to listen for notification actions: %v", err)
		os.Exit(1)
	}
}
