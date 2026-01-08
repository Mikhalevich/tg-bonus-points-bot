package event

type OutboxMessageType string

const (
	OutboxMessageTypePlain    OutboxMessageType = "plain"
	OutboxMessageTypeMarkdown OutboxMessageType = "markdown"
)

type OutboxMessage struct {
	ID             int               `json:"id"`
	ChatID         int64             `json:"chat_id"`
	ReplyMessageID *int64            `json:"reply_msg_id"`
	MessageText    string            `json:"msg_text"`
	MessageType    OutboxMessageType `json:"msg_type"`
	Payload        *string           `json:"payload"`
	Buttons        string            `json:"buttons"`
}
