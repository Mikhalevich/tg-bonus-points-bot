package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) SendPNG(
	ctx context.Context,
	chatID msginfo.ChatID,
	caption string,
	png []byte,
	rows ...button.ButtonRow,
) error {
	inlineButtons, err := m.SetButtonRows(ctx, rows...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if err := m.sender.SendPNGMarkdown(ctx, chatID, caption, png, inlineButtons...); err != nil {
		return fmt.Errorf("send png markdown: %w", err)
	}

	return nil
}
