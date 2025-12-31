package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) ReplyMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	textType MessageTextType,
	rows ...button.ButtonRow,
) error {
	inlineButtons, err := m.SetButtonRows(ctx, rows...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	switch textType {
	case MessageTextTypePlain:
		m.sender.ReplyText(ctx, chatID, replyMessageID, text, inlineButtons...)

	case MessageTextTypeMarkdown:
		m.sender.ReplyTextMarkdown(ctx, chatID, replyMessageID, text, inlineButtons...)
	}

	return nil
}
