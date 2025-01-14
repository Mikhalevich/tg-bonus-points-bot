package buttonrespository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

func (r *ButtonRepository) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) ([]button.InlineKeyboardButtonRow, error) {
	key := generateID()

	inlineButtons, hMap, err := processButtonRows(key, rows)
	if err != nil {
		return nil, fmt.Errorf("process button rows: %w", err)
	}

	if err := r.client.HSet(ctx, key, hMap).Err(); err != nil {
		return nil, fmt.Errorf("hset: %w", err)
	}

	return inlineButtons, nil
}

func processButtonRows(
	key string,
	rows []button.ButtonRow,
) ([]button.InlineKeyboardButtonRow, map[string]any, error) {
	var (
		inlineButtons = make([]button.InlineKeyboardButtonRow, 0, len(rows))
		hMap          = make(map[string]any)
		buttonNum     = int64(1)
	)

	for _, buttonsRow := range rows {
		inlineButtonsRow := make([]button.InlineKeyboardButton, 0, len(buttonsRow))

		for _, b := range buttonsRow {
			formattedNum := strconv.FormatInt(buttonNum, 10)

			inlineButtonsRow = append(inlineButtonsRow, button.InlineKeyboardButton{
				ID:      buttonID(key, formattedNum),
				Caption: b.Caption,
			})

			encodedButton, err := encodeButton(b)
			if err != nil {
				return nil, nil, fmt.Errorf("encode button: %w", err)
			}

			hMap[formattedNum] = encodedButton

			buttonNum++
		}

		inlineButtons = append(inlineButtons, inlineButtonsRow)
	}

	return inlineButtons, hMap, nil
}

func buttonID(key string, num string) button.ID {
	return button.IDFromString(fmt.Sprintf("%s_%s", key, num))
}
