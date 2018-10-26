package xiaomi

import (
	"encoding/json"
	"strconv"
	"strings"
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

func (m *Message) SetRestrictedPackageName(restrictedPackageNames []string) *Message {
	m.RestrictedPackageName = strings.Join(restrictedPackageNames, ",")
	return m
}

func (m *Message) SetPassThrough(passThrough int32) *Message {
	m.PassThrough = passThrough
	return m
}

func (m *Message) SetNotifyType(notifyType int32) *Message {
	m.NotifyType = notifyType
	return m
}

func (m *Message) SetTimeToSend(tts int64) *Message {
	if time.Since(time.Unix(0, tts*int64(time.Millisecond))) > MaxTimeToSend {
		m.TimeToSend = time.Now().Add(MaxTimeToSend).UnixNano() / 1e6
	} else {
		m.TimeToSend = tts
	}
	return m
}

func (m *Message) SetTimeToLive(ttl int64) *Message {
	if time.Since(time.Unix(0, ttl*int64(time.Millisecond))) > MaxTimeToLive {
		m.TimeToLive = time.Now().Add(MaxTimeToLive).UnixNano() / 1e6
	} else {
		m.TimeToLive = ttl
	}
	return m
}

func (m *Message) SetNotifyID(notifyID int64) *Message {
	m.NotifyID = notifyID
	return m
}

func (m *Message) EnableFlowControl() *Message {
	m.Extra["flow_control"] = "1"
	return m
}

func (m *Message) DisableFlowControl() *Message {
	delete(m.Extra, "flow_control")
	return m
}

func (m *Message) SetJobKey(jobKey string) *Message {
	m.Extra["jobkey"] = jobKey
	return m
}

func (m *Message) SetCallback(callbackURL string) *Message {
	m.Extra["callback"] = callbackURL
	m.Extra["callback.type"] = "3"
	return m
}

func (m *Message) AddExtra(key, value string) *Message {
	m.Extra[key] = value
	return m
}

func (m *Message) JSON() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes
}

func NewAndroidMessage(title, description string) *Message {
	return &Message{
		Payload:     "",
		Title:       title,
		Description: description,
		PassThrough: 0,
		NotifyType:  -1, // default notify type
		TimeToLive:  0,
		TimeToSend:  0,
		NotifyID:    0,
		Extra:       make(map[string]string),
	}
}

func (m *Message) SetLauncherActivity() *Message {
	m.Extra["notify_effect"] = "1"
	return m
}

func (m *Message) SetJumpActivity(value string) *Message {
	m.Extra["notify_effect"] = "2"
	m.Extra["intent_uri"] = value
	return m
}

func (m *Message) SetJumpWebURL(value string) *Message {
	m.Extra["notify_effect"] = "3"
	m.Extra["web_uri"] = value
	return m
}

func (m *Message) SetPayload(payload string) *Message {
	m.Payload = payload
	return m
}

func NewIOSMessage(description string) *Message {
	return &Message{
		Payload:     "",
		Title:       "",
		Description: description,
		PassThrough: 0,
		NotifyType:  -1, // default notify type
		TimeToLive:  0,
		TimeToSend:  0,
		NotifyID:    0,
		Extra:       make(map[string]string),
	}
}

func (i *Message) SetBadge(badge int64) *Message {
	i.Extra["badge"] = strconv.FormatInt(badge, 10)
	return i
}

func (i *Message) SetCategory(category string) *Message {
	i.Extra["category"] = category
	return i
}

func (i *Message) SetSoundURL(soundURL string) *Message {
	i.Extra["sound_url"] = soundURL
	return i
}

type TargetType int32

const (
	TargetTypeRegID   TargetType = 1
	TargetTypeReAlias TargetType = 2
	TargetTypeAccount TargetType = 3
)

type TargetedMessage struct {
	message    *Message
	targetType TargetType
	target     string
}

func NewTargetedMessage(m *Message, target string, targetType TargetType) *TargetedMessage {
	return &TargetedMessage{
		message:    m,
		targetType: targetType,
		target:     target,
	}
}

func (tm *TargetedMessage) SetTargetType(targetType TargetType) *TargetedMessage {
	tm.targetType = targetType
	return tm
}

func (tm *TargetedMessage) SetTarget(target string) *TargetedMessage {
	tm.target = target
	return tm
}

func (tm *TargetedMessage) JSON() []byte {
	bytes, err := json.Marshal(tm)
	if err != nil {
		panic(err)
	}
	return bytes
}
