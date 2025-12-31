package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) SendMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
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
		m.sender.SendText(ctx, chatID, text, inlineButtons...)

	case MessageTextTypeMarkdown:
		m.sender.SendTextMarkdown(ctx, chatID, text, inlineButtons...)
	}

	return nil
}
