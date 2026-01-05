package messagesender

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

var (
	_ messageprocessor.Sender          = (*messageSender)(nil)
	_ messageprocessor.MarkdownEscaper = (*messageSender)(nil)
)

type messageSender struct {
	bot          *bot.Bot
	paymentToken string
}

func New(bot *bot.Bot, paymentToken string) *messageSender {
	return &messageSender{
		bot:          bot,
		paymentToken: paymentToken,
	}
}

func (m *messageSender) ReplyText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_plain", text).
			Error("failed to reply text")
	}
}

func makeButtonsMarkup(rows ...button.InlineKeyboardButtonRow) models.ReplyMarkup {
	if len(rows) == 0 {
		return nil
	}

	keyboard := make([][]models.InlineKeyboardButton, 0, len(rows))

	for _, row := range rows {
		buttonRow := make([]models.InlineKeyboardButton, 0, len(row))

		for _, b := range row {
			buttonRow = append(buttonRow, models.InlineKeyboardButton{
				Text:         b.Caption,
				CallbackData: b.ID.String(),
				Pay:          b.Pay,
			})
		}

		keyboard = append(keyboard, buttonRow)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func (m *messageSender) ReplyTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
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
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID.Int64(),
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
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
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID.Int64(),
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
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
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: chatID.Int64(),
		Photo: &models.InputFileUpload{
			Data: bytes.NewReader(png),
		},
		Caption:     caption,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send photo: %w", err)
	}

	return nil
}
