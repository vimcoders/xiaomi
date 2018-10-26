package xiaomi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type MiPush struct {
	packageName []string
	host        string
	appSecret   string
}

func NewClient(appSecret string, packageName []string) *MiPush {
	return &MiPush{
		packageName: packageName,
		host:        ProductionHost,
		appSecret:   appSecret,
	}
}

func (m *MiPush) Send(ctx context.Context, msg *Message, regID string) (*SendResult, error) {
	params := m.assembleSendParams(msg, regID)
	bytes, err := m.doPost(ctx, m.host+RegURL, params)
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) doPost(ctx context.Context, url string, form url.Values) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)
	client := &http.Client{}
	tryTime := 0
tryAgain:
	res, err = ctxhttp.Do(ctx, client, req)
	if err != nil {
		fmt.Println("xiaomi push post err:", err, tryTime)
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}
		tryTime += 1
		if tryTime < PostRetryTimes {
			goto tryAgain
		}
		return nil, err
	}
	if res.Body == nil {
		panic("xiaomi response is nil")
	}
	defer res.Body.Close()
	fmt.Println("res.StatusCode=", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) doGet(ctx context.Context, url string, params string) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("GET", url+params, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)

	client := &http.Client{}
	res, err = ctxhttp.Do(ctx, client, req)
	if res.Body == nil {
		panic("xiaomi response is nil")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) assembleSendParams(msg *Message, regID string) url.Values {
	form := m.defaultForm(msg)
	form.Add("registration_id", regID)
	return form
}

func (m *MiPush) defaultForm(msg *Message) url.Values {
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
	return form
}
