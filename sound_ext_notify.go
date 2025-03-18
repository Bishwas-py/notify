// Package notify Extends the notify package to include sound notifications
package notify

import (
	"github.com/godbus/dbus/v5"
	"log"
)

type SoundName string

// Standard freedesktop sound names
const (
	AlarmClockElapsed              SoundName = "alarm-clock-elapsed"
	AudioChannelFrontCenter        SoundName = "audio-channel-front-center"
	AudioChannelFrontLeft          SoundName = "audio-channel-front-left"
	AudioChannelFrontRight         SoundName = "audio-channel-front-right"
	AudioChannelRearCenter         SoundName = "audio-channel-rear-center"
	AudioChannelRearLeft           SoundName = "audio-channel-rear-left"
	AudioChannelRearRight          SoundName = "audio-channel-rear-right"
	AudioChannelSideLeft           SoundName = "audio-channel-side-left"
	AudioChannelSideRight          SoundName = "audio-channel-side-right"
	AudioTestSignal                SoundName = "audio-test-signal"
	AudioVolumeChange              SoundName = "audio-volume-change"
	Bell                           SoundName = "bell"
	CameraShutter                  SoundName = "camera-shutter"
	Complete                       SoundName = "complete"
	DeviceAdded                    SoundName = "device-added"
	DeviceRemoved                  SoundName = "device-removed"
	DialogError                    SoundName = "dialog-error"
	DialogInformation              SoundName = "dialog-information"
	DialogWarning                  SoundName = "dialog-warning"
	Message                        SoundName = "message"
	MessageNewInstant              SoundName = "message-new-instant"
	NetworkConnectivityEstablished SoundName = "network-connectivity-established"
	NetworkConnectivityLost        SoundName = "network-connectivity-lost"
	PhoneIncomingCall              SoundName = "phone-incoming-call"
	PhoneOutgoingBusy              SoundName = "phone-outgoing-busy"
	PhoneOutgoingCalling           SoundName = "phone-outgoing-calling"
	PowerPlug                      SoundName = "power-plug"
	PowerUnplug                    SoundName = "power-unplug"
	ScreenCapture                  SoundName = "screen-capture"
	ServiceLogin                   SoundName = "service-login"
	ServiceLogout                  SoundName = "service-logout"
	SuspendError                   SoundName = "suspend-error"
	TrashEmpty                     SoundName = "trash-empty"
	WindowAttention                SoundName = "window-attention"
	WindowQuestion                 SoundName = "window-question"
)

func (n *Notification) SetSoundByName(name SoundName) {
	if n.Hints == nil {
		n.Hints = make(map[string]dbus.Variant)
	}
	n.Hints["sound-name"] = dbus.MakeVariant(string(name))
	log.Printf("%v", n.Hints)
}

func (n *Notification) SetSoundByPath(soundFilePath string) {
	if n.Hints == nil {
		n.Hints = make(map[string]dbus.Variant)
	}
	n.Hints["sound-file"] = dbus.MakeVariant(soundFilePath)
}
