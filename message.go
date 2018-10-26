package xiaomi

import (
	"time"
)

type Message struct {
	RestrictedPackageName string            `json:"restricted_package_name,omitempty"`
	Payload               string            `json:"payload,omitempty"`
	Title                 string            `json:"title,omitempty"`
	Description           string            `json:"description,omitempty"`
	PassThrough           int32             `json:"pass_through"`
	NotifyType            int32             `json:"notify_type,omitempty"`
	TimeToLive            int64             `json:"time_to_live,omitempty"`
	TimeToSend            int64             `json:"time_to_send,omitempty"`
	NotifyID              int64             `json:"notify_id"`
	Extra                 map[string]string `json:"extra,omitempty"`
}

const (
	MaxTimeToSend = time.Hour * 24 * 7
	MaxTimeToLive = time.Hour * 24 * 7 * 2
)

func NewAndroidMessage(title, description string) *Message {
	return &Message{
		Payload:     "",
		Title:       title,
		Description: description,
		PassThrough: 0,
		NotifyType:  -1,
		TimeToLive:  0,
		TimeToSend:  0,
		NotifyID:    0,
		Extra:       make(map[string]string),
	}
}

func NewIOSMessage(description string) *Message {
	return &Message{
		Payload:     "",
		Title:       "",
		Description: description,
		PassThrough: 0,
		NotifyType:  -1,
		TimeToLive:  0,
		TimeToSend:  0,
		NotifyID:    0,
		Extra:       make(map[string]string),
	}
}
