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

	if err := m.sendMsg(ctx, chatID, text, textType, inlineButtons); err != nil {
		return fmt.Errorf("send msg: %w", err)
	}

	return nil
}

func (m *MessageProcessor) sendMsg(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	textType MessageTextType,
	inlineButtons []button.InlineKeyboardButtonRow,
) error {
	switch textType {
	case MessageTextTypePlain:
		if err := m.sender.SendText(ctx, chatID, text, inlineButtons...); err != nil {
			return fmt.Errorf("send text: %w", err)
		}

	case MessageTextTypeMarkdown:
		if err := m.sender.SendTextMarkdown(ctx, chatID, text, inlineButtons...); err != nil {
			return fmt.Errorf("send text markdown: %w", err)
		}
	}

	return nil
}
