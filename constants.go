package xiaomi

const (
	ProductionHost = "https://api.xmpush.xiaomi.com"
)

const (
	RegURL                               = "/v3/message/regid"
	MultiMessagesRegIDURL                = "/v2/multi_messages/regids"
	MultiMessagesAliasURL                = "/v2/multi_messages/aliases"
	MultiMessagesUserAccountURL          = "/v2/multi_messages/user_accounts"
	MessageAlisaURL                      = "/v3/message/alias"
	MessageUserAccountURL                = "/v2/message/user_account"
	MultiPackageNameMessageMultiTopicURL = "/v3/message/multi_topic"
	MessageMultiTopicURL                 = "/v2/message/topic"
	MultiPackageNameMessageAllURL        = "/v3/message/all"
	MessageAllURL                        = "/v2/message/all"
	MultiTopicURL                        = "/v3/message/multi_topic"
	ScheduleJobExistURL                  = "/v2/schedule_job/exist"
	ScheduleJobDeleteURL                 = "/v2/schedule_job/delete"
	ScheduleJobDeleteByJobKeyURL         = "/v3/schedule_job/delete"
)

const (
	StatsURL          = "/v1/stats/message/counters"
	MessageStatusURL  = "/v1/trace/message/status"
	MessagesStatusURL = "/v1/trace/messages/status"
)

const (
	TopicSubscribeURL          = "/v2/topic/subscribe"
	TopicUnSubscribeURL        = "/v2/topic/unsubscribe"
	TopicSubscribeByAliasURL   = "/v2/topic/subscribe/alias"
	TopicUnSubscribeByAliasURL = "/v2/topic/unsubscribe/alias"
)

const (
	InvalidRegIDsURL = "https://feedback.xmpush.xiaomi.com/v1/feedback/fetch_invalid_regids"
)

const (
	AliasAllURL  = "/v1/alias/all"
	TopicsAllURL = "/v1/topic/all"
)

var (
	PostRetryTimes = 3
)

var (
	BrandsMap = map[string]string{
		"品牌":    "MODEL",
		"小米":    "xiaomi",
		"三星":    "samsung",
		"华为":    "huawei",
		"中兴":    "zte",
		"中兴努比亚": "nubia",
		"酷派":    "coolpad",
		"联想":    "lenovo",
		"魅族":    "meizu",
		"HTC":   "htc",
		"OPPO":  "oppo",
		"VIVO":  "vivo",
		"摩托罗拉":  "motorola",
		"索尼":    "sony",
		"LG":    "lg",
		"金立":    "jinli",
		"天语":    "tianyu",
		"诺基亚":   "nokia",
		"美图秀秀":  "meitu",
		"谷歌":    "google",
		"TCL":   "tcl",
		"锤子手机":  "chuizi",
		"一加手机":  "1+",
		"中国移动":  "chinamobile",
		"昂达":    "angda",
		"邦华":    "banghua",
		"波导":    "bird",
		"长虹":    "changhong",
		"大可乐":   "dakele",
		"朵唯":    "doov",
		"海尔":    "haier",
		"海信":    "hisense",
		"康佳":    "konka",
		"酷比魔方":  "kubimofang",
		"米歌":    "mige",
		"欧博信":   "ouboxin",
		"欧新":    "ouxin",
		"飞利浦":   "philip",
		"维图":    "voto",
		"小辣椒":   "xiaolajiao",
		"夏新":    "xiaxin",
		"亿通":    "yitong",
		"语信":    "yuxin",
	}

	PriceMap = map[string]string{
		"0-999":     "0-999",
		"1000-1999": "1000-1999",
		"2000-3999": "2000-3999",
		"4000+":     "4000+",
	}
)
