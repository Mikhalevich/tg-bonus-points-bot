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

	if err := m.replyMsg(ctx, chatID, replyMessageID, text, textType, inlineButtons); err != nil {
		return fmt.Errorf("reply msg: %w", err)
	}

	return nil
}

func (m *MessageProcessor) replyMsg(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	textType MessageTextType,
	inlineButtons []button.InlineKeyboardButtonRow,
) error {
	switch textType {
	case MessageTextTypePlain:
		if err := m.sender.ReplyText(
			ctx,
			chatID,
			replyMessageID,
			text,
			inlineButtons...,
		); err != nil {
			return fmt.Errorf("reply text: %w", err)
		}

	case MessageTextTypeMarkdown:
		if err := m.sender.ReplyTextMarkdown(
			ctx,
			chatID,
			replyMessageID,
			text,
			inlineButtons...,
		); err != nil {
			return fmt.Errorf("reply text markdown: %w", err)
		}
	}

	return nil
}
