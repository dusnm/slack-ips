package slack

import "net/url"

type (
	Message struct {
		Token          string
		TeamID         string
		TeamDomain     string
		EnterpriseID   string
		EnterpriseName string
		ChannelID      string
		ChannelName    string
		UserID         string
		UserName       string
		Command        string
		Text           string
		ResponseURL    string
		TriggerID      string
		APIAppID       string
	}
)

func NewMessage(v url.Values) Message {
	return Message{
		Token:          v.Get("token"),
		TeamID:         v.Get("team_id"),
		TeamDomain:     v.Get("team_domain"),
		EnterpriseID:   v.Get("enterprise_id"),
		EnterpriseName: v.Get("enterprise_name"),
		ChannelID:      v.Get("channel_id"),
		ChannelName:    v.Get("channel_name"),
		UserID:         v.Get("user_id"),
		UserName:       v.Get("user_name"),
		Command:        v.Get("command"),
		Text:           v.Get("text"),
		ResponseURL:    v.Get("response_url"),
		TriggerID:      v.Get("trigger_id"),
		APIAppID:       v.Get("api_app_id"),
	}
}
