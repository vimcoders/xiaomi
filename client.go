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

func (m *MiPush) SendToList(ctx context.Context, msg *Message, regIDList []string) (*SendResult, error) {
	if len(regIDList) == 0 || len(regIDList) > 1000 {
		panic("wrong number regIDList")
	}
	return m.Send(ctx, msg, strings.Join(regIDList, ","))
}

func (m *MiPush) SendTargetMessageList(ctx context.Context, msgList []*TargetedMessage) (*SendResult, error) {
	if len(msgList) == 0 {
		return nil, errors.New("empty msg")
	}
	if len(msgList) == 1 {
		return m.Send(ctx, msgList[0].message, msgList[0].target)
	}
	params := m.assembleTargetMessageListParams(msgList)
	var bytes []byte
	var err error
	if msgList[0].targetType == TargetTypeRegID {
		bytes, err = m.doPost(ctx, m.host+MultiMessagesRegIDURL, params)
	} else if msgList[0].targetType == TargetTypeReAlias {
		bytes, err = m.doPost(ctx, m.host+MultiMessagesAliasURL, params)
	} else if msgList[0].targetType == TargetTypeAccount {
		bytes, err = m.doPost(ctx, m.host+MultiMessagesUserAccountURL, params)
	} else {
		panic("bad targetType")
	}

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

func (m *MiPush) SendToAlias(ctx context.Context, msg *Message, alias string) (*SendResult, error) {
	params := m.assembleSendToAlisaParams(msg, alias)
	bytes, err := m.doPost(ctx, m.host+MessageAlisaURL, params)
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

func (m *MiPush) SendToAliasList(ctx context.Context, msg *Message, aliasList []string) (*SendResult, error) {
	if len(aliasList) == 0 || len(aliasList) > 1000 {
		panic("wrong number aliasList")
	}
	return m.SendToAlias(ctx, msg, strings.Join(aliasList, ","))
}

func (m *MiPush) SendToUserAccount(ctx context.Context, msg *Message, userAccount string) (*SendResult, error) {
	params := m.assembleSendToUserAccountParams(msg, userAccount)
	bytes, err := m.doPost(ctx, m.host+MessageUserAccountURL, params)
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

func (m *MiPush) SendToUserAccountList(ctx context.Context, msg *Message, accountList []string) (*SendResult, error) {
	if len(accountList) == 0 || len(accountList) > 1000 {
		panic("wrong number accountList")
	}
	return m.SendToUserAccount(ctx, msg, strings.Join(accountList, ","))
}

func (m *MiPush) Broadcast(ctx context.Context, msg *Message, topic string) (*SendResult, error) {
	params := m.assembleBroadcastParams(msg, topic)
	var bytes []byte
	var err error
	if len(m.packageName) > 1 {
		bytes, err = m.doPost(ctx, m.host+MultiPackageNameMessageMultiTopicURL, params)
	} else {
		bytes, err = m.doPost(ctx, m.host+MessageMultiTopicURL, params)
	}
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

func (m *MiPush) BroadcastAll(ctx context.Context, msg *Message) (*SendResult, error) {
	params := m.assembleBroadcastAllParams(msg)
	var bytes []byte
	var err error
	if len(m.packageName) > 1 {
		bytes, err = m.doPost(ctx, m.host+MultiPackageNameMessageAllURL, params)
	} else {
		bytes, err = m.doPost(ctx, m.host+MessageAllURL, params)
	}
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

type TopicOP string

const (
	UNION        TopicOP = "UNION"
	INTERSECTION TopicOP = "INTERSECTION"
	EXCEPT       TopicOP = "EXCEPT"
)

func (m *MiPush) MultiTopicBroadcast(ctx context.Context, msg *Message, topics []string, topicOP TopicOP) (*SendResult, error) {
	if len(topics) > 5 || len(topics) == 0 {
		panic("topics size invalid")
	}
	if len(topics) == 1 {
		return m.Broadcast(ctx, msg, topics[0])
	}
	params := m.assembleMultiTopicBroadcastParams(msg, topics, topicOP)
	bytes, err := m.doPost(ctx, m.host+MultiTopicURL, params)
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

func (m *MiPush) CheckScheduleJobExist(ctx context.Context, msgID string) (*Result, error) {
	params := m.assembleCheckScheduleJobParams(msgID)
	bytes, err := m.doPost(ctx, m.host+ScheduleJobExistURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) DeleteScheduleJob(ctx context.Context, msgID string) (*Result, error) {
	params := m.assembleDeleteScheduleJobParams(msgID)
	bytes, err := m.doPost(ctx, m.host+ScheduleJobDeleteURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) DeleteScheduleJobByJobKey(ctx context.Context, jobKey string) (*Result, error) {
	params := m.assembleDeleteScheduleJobByJobKeyParams(jobKey)
	bytes, err := m.doPost(ctx, m.host+ScheduleJobDeleteByJobKeyURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) Stats(ctx context.Context, start, end, packageName string) (*StatsResult, error) {
	params := m.assembleStatsParams(start, end, packageName)
	bytes, err := m.doGet(ctx, m.host+StatsURL, params)
	if err != nil {
		return nil, err
	}
	var result StatsResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) GetMessageStatusByMsgID(ctx context.Context, msgID string) (*SingleStatusResult, error) {
	params := m.assembleStatusParams(msgID)
	bytes, err := m.doGet(ctx, m.host+MessageStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result SingleStatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) GetMessageStatusByJobKey(ctx context.Context, jobKey string) (*BatchStatusResult, error) {
	params := m.assembleStatusByJobKeyParams(jobKey)
	bytes, err := m.doGet(ctx, m.host+MessagesStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result BatchStatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) GetMessageStatusPeriod(ctx context.Context, beginTime, endTime int64) (*BatchStatusResult, error) {
	params := m.assembleStatusPeriodParams(beginTime, endTime)
	bytes, err := m.doGet(ctx, m.host+MessagesStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result BatchStatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) SubscribeTopicForRegID(ctx context.Context, regID, topic, category string) (*Result, error) {
	params := m.assembleSubscribeTopicForRegIDParams(regID, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicSubscribeURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) SubscribeTopicForRegIDList(ctx context.Context, regIDList []string, topic, category string) (*Result, error) {
	return m.SubscribeTopicForRegID(ctx, strings.Join(regIDList, ","), topic, category)
}

func (m *MiPush) UnSubscribeTopicForRegID(ctx context.Context, regID, topic, category string) (*Result, error) {
	params := m.assembleUnSubscribeTopicForRegIDParams(regID, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicUnSubscribeURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) UnSubscribeTopicForRegIDList(ctx context.Context, regIDList []string, topic, category string) (*Result, error) {
	return m.UnSubscribeTopicForRegID(ctx, strings.Join(regIDList, ","), topic, category)
}

func (m *MiPush) SubscribeTopicByAlias(ctx context.Context, aliases []string, topic, category string) (*Result, error) {
	params := m.assembleSubscribeTopicByAliasParams(aliases, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicSubscribeByAliasURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) UnSubscribeTopicByAlias(ctx context.Context, aliases []string, topic, category string) (*Result, error) {
	params := m.assembleUnSubscribeTopicByAliasParams(aliases, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicUnSubscribeByAliasURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) GetInvalidRegIDs(ctx context.Context) (*InvalidRegIDsResult, error) {
	params := m.assembleGetInvalidRegIDsParams()
	bytes, err := m.doGet(ctx, InvalidRegIDsURL, params)
	if err != nil {
		return nil, err
	}
	var result InvalidRegIDsResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) GetAliasesOfRegID(ctx context.Context, regID string) (*AliasesOfRegIDResult, error) {
	params := m.assembleGetAliasesOfParams(regID)
	bytes, err := m.doGet(ctx, m.host+AliasAllURL, params)
	if err != nil {
		return nil, err
	}
	var result AliasesOfRegIDResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) GetTopicsOfRegID(ctx context.Context, regID string) (*TopicsOfRegIDResult, error) {
	params := m.assembleGetTopicsOfParams(regID)
	bytes, err := m.doGet(ctx, m.host+TopicsAllURL, params)
	if err != nil {
		return nil, err
	}
	var result TopicsOfRegIDResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) assembleSendParams(msg *Message, regID string) url.Values {
	form := m.defaultForm(msg)
	form.Add("registration_id", regID)
	return form
}

func (m *MiPush) assembleTargetMessageListParams(msgList []*TargetedMessage) url.Values {
	form := url.Values{}
	type OneMsg struct {
		Target  string   `json:"target"`
		Message *Message `json:"message"`
	}
	var messages []*OneMsg

	for _, m := range msgList {
		messages = append(messages, &OneMsg{
			Target:  m.target,
			Message: m.message,
		})
	}
	bytes, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}
	form.Add("messages", string(bytes))
	form.Add("time_to_send", strconv.FormatInt(msgList[0].message.TimeToSend, 10))
	return form
}

func (m *MiPush) assembleSendToAlisaParams(msg *Message, alias string) url.Values {
	form := m.defaultForm(msg)
	form.Add("alias", alias)
	return form
}

func (m *MiPush) assembleSendToUserAccountParams(msg *Message, userAccount string) url.Values {
	form := m.defaultForm(msg)
	form.Add("user_account", userAccount)
	return form
}

func (m *MiPush) assembleBroadcastParams(msg *Message, topic string) url.Values {
	form := m.defaultForm(msg)
	form.Add("topic", topic)
	return form
}

func (m *MiPush) assembleBroadcastAllParams(msg *Message) url.Values {
	form := m.defaultForm(msg)
	return form
}

func (m *MiPush) assembleMultiTopicBroadcastParams(msg *Message, topics []string, topicOP TopicOP) url.Values {
	form := m.defaultForm(msg)
	form.Add("topic_op", string(topicOP))
	form.Add("topics", strings.Join(topics, ";$;"))
	return form
}

func (m *MiPush) assembleCheckScheduleJobParams(msgID string) url.Values {
	form := url.Values{}
	form.Add("job_id", msgID)
	return form
}

func (m *MiPush) assembleDeleteScheduleJobParams(msgID string) url.Values {
	form := url.Values{}
	form.Add("job_id", msgID)
	return form
}

func (m *MiPush) assembleDeleteScheduleJobByJobKeyParams(jobKey string) url.Values {
	form := url.Values{}
	form.Add("job_key", jobKey)
	return form
}

func (m *MiPush) assembleStatsParams(start, end, packageName string) string {
	form := url.Values{}
	form.Add("start_date", start)
	form.Add("end_date", end)
	form.Add("restricted_package_name", packageName)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusParams(msgID string) string {
	form := url.Values{}
	form.Add("msg_id", msgID)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusByJobKeyParams(jobKey string) string {
	form := url.Values{}
	form.Add("job_key", jobKey)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusPeriodParams(beginTime, endTime int64) string {
	form := url.Values{}
	form.Add("begin_time", strconv.FormatInt(int64(beginTime), 10))
	form.Add("end_time", strconv.FormatInt(int64(endTime), 10))
	return "?" + form.Encode()
}

func (m *MiPush) assembleSubscribeTopicForRegIDParams(regID, topic, category string) url.Values {
	form := url.Values{}
	form.Add("registration_id", regID)
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleUnSubscribeTopicForRegIDParams(regID, topic, category string) url.Values {
	form := url.Values{}
	form.Add("registration_id", regID)
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleSubscribeTopicByAliasParams(aliases []string, topic, category string) url.Values {
	form := url.Values{}
	form.Add("aliases", strings.Join(aliases, ","))
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleUnSubscribeTopicByAliasParams(aliases []string, topic, category string) url.Values {
	form := url.Values{}
	form.Add("aliases", strings.Join(aliases, ","))
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleGetInvalidRegIDsParams() string {
	form := url.Values{}
	return "?" + form.Encode()
}

func (m *MiPush) assembleGetAliasesOfParams(regID string) string {
	form := url.Values{}
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	form.Add("registration_id", regID)
	return "?" + form.Encode()
}

func (m *MiPush) assembleGetTopicsOfParams(regID string) string {
	form := url.Values{}
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	form.Add("registration_id", regID)
	return "?" + form.Encode()
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
