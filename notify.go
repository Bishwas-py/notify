package notify

import "github.com/godbus/dbus/v5"

type Notification struct {
	AppID     string
	ReplaceID string
	AppIcon   string
	Title     string
	Body      string
	Actions   Actions
	Hints     map[string]dbus.Variant
	Timeout   int
}

func (n *Notification) Trigger() (uint32, error) {
	listenHandlers, dbusActions := n.Actions.Results()
	notificationID, _ := notificationHandler.SendNotification(
		n.AppID,
		n.Title,
		n.Body,
		n.AppIcon,
		dbusActions,
		n.Hints,
		n.Timeout,
	)

	if err := notificationHandler.Listen(listenHandlers, notificationID, n.Timeout); err != nil {
		return 0, err
	}
	return notificationID, nil
}
