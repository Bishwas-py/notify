package main

import (
	"fmt"
	"github.com/Bishwas-py/notify"
	"github.com/godbus/dbus/v5"
	"log"
	"os/exec"
	"time"
)

func LogoutViaGnome() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %v", err)
	}

	obj := conn.Object("org.gnome.SessionManager", "/org/gnome/SessionManager")
	call := obj.Call("org.gnome.SessionManager.Logout", 0, uint32(0))

	return call.Err
}

func LogoutViaSystemd() error {
	cmd := exec.Command("loginctl", "terminate-user", "$USER")
	return cmd.Run()
}

func HandleLogout() {
	log.Println("Logout action triggered!")

	err := LogoutViaGnome()
	if err != nil {
		log.Printf("GNOME logout failed: %v\n", err)

		err = LogoutViaSystemd()
		if err != nil {
			log.Printf("Systemd logout failed: %v\n", err)
		}
	}
}

func HandleCall() {
	log.Println("Call action triggered!")
}

func main() {
	actions := notify.Actions{
		{
			Title:   "Logout Now",
			Trigger: HandleLogout,
		},
		{
			Trigger: HandleCall,
			Title:   "Call Me",
		},
	}

	timeout := 10 * time.Second

	notification := notify.Notification{
		AppID:   "example",
		Title:   "Hello, World!",
		Body:    "This is an example notification",
		Timeout: int(timeout),
		Actions: actions,
	}
	notification.SetSoundByName("alarm-clock-elapsed")

	_, _ = notification.Trigger()
}
