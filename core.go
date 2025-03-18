package notify

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"log"
	"reflect"
	"runtime"
	"time"
)

type NotificationHandler struct {
	conn *dbus.Conn
}

var notificationHandler *NotificationHandler

func init() {
	log.Printf("Running on Go %s", runtime.Version())
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
func (h *NotificationHandler) SendNotification(
	appName string,
	title string,
	body string,
	appIcon string,
	dbusActions []string,
	variants map[string]dbus.Variant,
	timeout int,
) (uint32, error) {
	obj := h.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var id uint32
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		appName,     // app_name
		uint32(0),   // replaces_id
		appIcon,     // app_icon
		title,       // summary
		body,        // body
		dbusActions, // actions
		variants,    // hints
		timeout,
	)

	if err := call.Store(&id); err != nil {
		return 0, fmt.Errorf("failed to send notification: %v", err)
	}

	return id, nil
}

func (h *NotificationHandler) CloseNotification(id uint32) error {
	obj := h.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call("org.freedesktop.Notifications.CloseNotification", 0, id)
	return call.Err
}

// Listen starts listening for notification action signals
func (h *NotificationHandler) Listen(actionHandlers map[string]func(), notificationID uint32, timeout int) error {
	if err := h.conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.Notifications"),
		dbus.WithMatchMember("ActionInvoked"),
	); err != nil {
		return fmt.Errorf("failed to add match: %v", err)
	}

	c := make(chan *dbus.Signal, 1)
	h.conn.Signal(c)

	defer func() {
		_ = h.conn.RemoveMatchSignal(
			dbus.WithMatchInterface("org.freedesktop.Notifications"),
			dbus.WithMatchMember("ActionInvoked"),
		)
		h.conn.RemoveSignal(c)
	}()

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

			curNid := v.Body[0].(uint32)
			if curNid != notificationID {
				continue
			}
			if v.Name == "org.freedesktop.Notifications.ActionInvoked" {
				actionKey := v.Body[1].(string)

				log.Printf("Notification %d: Action %s invoked %d\n", curNid, actionKey, notificationID)
				if handler, ok := actionHandlers[actionKey]; ok {
					handler()
					_ = h.conn.Close()
					_ = h.CloseNotification(notificationID)
					break actionLoop
				}
			}
			if v.Name == "org.freedesktop.Notifications.NotificationClosed" {
				_ = h.CloseNotification(notificationID)
				break actionLoop
			}

		case <-timer.C:
			println("Timeout reached, breaking loop")
			_ = h.CloseNotification(notificationID)
			break actionLoop
		}
	}

	return nil
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

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
