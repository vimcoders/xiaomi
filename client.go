package xiaomi

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	POST_RETRY_TIMES = 3
	XIAOMI_REGID_URL = "/v3/message/regid"
	XIAOMI_HOST_URL  = "https://api.xmpush.xiaomi.com"
)

type MiPush struct {
	packageName []string
	host        string
	appSecret   string
	token       string
}

func NewClient(appSecret string, packageName []string, token string) *MiPush {
	return &MiPush{
		packageName: packageName,
		host:        XIAOMI_REGID_URL,
		appSecret:   appSecret,
		token:       token,
	}
}

func (m *MiPush) Send(ctx context.Context, msg *Message) (string, error) {
	for i := 0; i < POST_RETRY_TIMES; i++ {
		bytes, err := m.doPost(ctx, m.host+XIAOMI_REGID_URL, m.ToFormValues(msg))
		if err != nil {
			if i+1 == POST_RETRY_TIMES {
				return "", err
			}
			continue
		}
		return string(bytes), nil
	}

	return "", nil
}

func (m *MiPush) doPost(ctx context.Context, url string, form url.Values) (bytes []byte, result error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)
	client := &http.Client{}
	res, err := ctxhttp.Do(ctx, client, req)
	if err != nil {
		select {
		case <-ctx.Done():
		default:
		}
	}
	if res.Body == nil {
		panic("xiaomi response is nil")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (m *MiPush) ToFormValues(msg *Message) url.Values {
	form := url.Values{}
	if len(m.packageName) > 0 {
		form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	}
	if msg.TimeToLive > 0 {
		form.Add("time_to_live", strconv.FormatInt(msg.TimeToLive, 10))
	}
	if len(msg.Payload) > 0 {
		form.Add("payload", msg.Payload)
	}
	if len(msg.Title) > 0 {
		form.Add("title", msg.Title)
	}
	if len(msg.Description) > 0 {
		form.Add("description", msg.Description)
	}
	form.Add("notify_type", strconv.FormatInt(int64(msg.NotifyType), 10))
	form.Add("pass_through", strconv.FormatInt(int64(msg.PassThrough), 10))
	if msg.NotifyID != 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.NotifyID), 10))
	}
	if msg.TimeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.TimeToSend), 10))
	}
	if len(msg.Extra) > 0 {
		for k, v := range msg.Extra {
			form.Add("extra."+k, v)
		}
	}
	if len(m.token) > 0 {
		form.Add("registration_id", m.token)
	}
	return form
}
