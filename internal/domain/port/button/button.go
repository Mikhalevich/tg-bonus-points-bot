package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(s string) ID {
	return ID(s)
}

type Button struct {
	ChatID    msginfo.ChatID
	Caption   string
	Operation Operation
	Payload   []byte
}

type ButtonRow []Button

func Row(buttons ...Button) ButtonRow {
	return buttons
}

func GetPayload[P any](b Button) (P, error) {
	payload, err := gobDecodePayload[P](b.Payload)
	if err != nil {
		return payload, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}

func gobEncodePayload(p any) ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func gobDecodePayload[Payload any](b []byte) (Payload, error) {
	var payload Payload
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return payload, fmt.Errorf("gob decode: %w", err)
	}

	return payload, nil
}

func createButton[P any](
	chatID msginfo.ChatID,
	caption string,
	operation Operation,
	payload P,
) (Button, error) {
	payloadBytes, err := gobEncodePayload(payload)

	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: operation,
		Payload:   payloadBytes,
	}, nil
}
