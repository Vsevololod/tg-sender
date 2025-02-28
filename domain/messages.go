package domain

import (
	"context"
	msg "github.com/Vsevololod/tg-api-contracts-lib/gen/go/messages"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	//"google.golang.org/protobuf/proto"
)

//type Message struct {
//	Text   string            `json:"text"`
//	UserId uint64            `json:"user_id"`
//	Type   string            `json:"type"`
//	Params map[string]string `json:"params"`
//}

type MessageWithContext struct {
	Message *msg.TgSendMessage
	UUID    string
	Context context.Context
}

func ParseMessage(jsonData []byte, isJson bool) (*msg.TgSendMessage, error) {
	message := msg.TgSendMessage{}
	if isJson {
		err := protojson.Unmarshal(jsonData, &message)
		if err != nil {
			return nil, err
		}
	} else {
		err := proto.Unmarshal(jsonData, &message)
		if err != nil {
			return nil, err
		}
	}
	return &message, nil
}
