package messagesender

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

var _ port.MessageSender = (*messageSender)(nil)

type messageSender struct {
	bot *bot.Bot
}

func New(bot *bot.Bot) *messageSender {
	return &messageSender{
		bot: bot,
	}
}

func (m *messageSender) ReplyText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	buttons ...port.Button,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(buttons...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_plain", text).
			Error("failed to reply text")
	}
}

func makeButtonsMarkup(buttons ...port.Button) models.ReplyMarkup {
	if len(buttons) == 0 {
		return nil
	}

	buttonRow := make([]models.InlineKeyboardButton, 0, len(buttons))
	for _, b := range buttons {
		buttonRow = append(buttonRow, models.InlineKeyboardButton{
			Text:         b.Text,
			CallbackData: b.ID.String(),
		})
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			buttonRow,
		},
	}
}

func (m *messageSender) ReplyTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	buttons ...port.Button,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(buttons...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_markdown", text).
			Error("failed to reply text markdown")
	}
}

func (m *messageSender) SendTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID.Int64(),
		ParseMode: models.ParseModeMarkdown,
		Text:      text,
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_markdown", text).
			Error("failed to send text markdown")
	}
}

func (m *messageSender) SendText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		Text:   text,
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text", text).
			Error("failed to send text")
	}
}

func (m *messageSender) EscapeMarkdown(s string) string {
	return bot.EscapeMarkdown(s)
}

func (m *messageSender) SendPNGMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	caption string,
	png []byte,
	buttons ...port.Button,
) error {
	if _, err := m.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: chatID.Int64(),
		Photo: &models.InputFileUpload{
			Data: bytes.NewReader(png),
		},
		Caption:     caption,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: makeButtonsMarkup(buttons...),
	}); err != nil {
		return fmt.Errorf("send photo: %w", err)
	}

	return nil
}
