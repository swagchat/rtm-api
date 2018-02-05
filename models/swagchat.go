package models

import "github.com/swagchat/rtm-api/utils"

type Message struct {
	MessageId string         `json:"messageId,omitempty"`
	RoomId    string         `json:"roomId,omitempty"`
	UserId    string         `json:"userId,omitempty"`
	Type      string         `json:"type,omitempty"`
	EventName string         `json:"eventName,omitempty"`
	Payload   utils.JSONText `json:"payload,omitempty"`
	Created   int64          `json:"created,omitempty"`
	Modified  int64          `json:"modified,omitempty"`
}

type PayloadText struct {
	Text string `json:"text"`
}

type PayloadImage struct {
	Mime         string `json:"mime"`
	SourceUrl    string `json:"sourceUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type PayloadUsers struct {
	Users []UserForRoom `json:"users"`
}

type UserForRoom struct {
	// from User
	UserId         string         `json:"userId,omitempty"`
	Name           string         `json:"name,omitempty"`
	PictureUrl     string         `json:"pictureUrl,omitempty"`
	InformationUrl string         `json:"informationUrl,omitempty"`
	MetaData       utils.JSONText `json:"metaData,omitempty"`
	Created        int64          `json:"created"`
	Modified       int64          `json:"modified,omitempty"`

	// from RoomUser
	RuUnreadCount *int64         `json:"ruUnreadCount,omitempty"`
	RuMetaData    utils.JSONText `json:"ruMetaData,omitempty"`
	RuCreated     int64          `json:"ruCreated,omitempty"`
	RuModified    int64          `json:"ruModified,omitempty"`
}
