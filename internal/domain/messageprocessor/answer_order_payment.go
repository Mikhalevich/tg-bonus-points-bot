package messageprocessor

import (
	"context"
	"fmt"
)

func (m *MessageProcessor) AnswerOrderPayment(
	ctx context.Context,
	paymentID string,
	ok bool,
	errorMsg string,
) error {
	if err := m.sender.AnswerOrderPayment(ctx, paymentID, ok, errorMsg); err != nil {
		return fmt.Errorf("answer order payment: %w", err)
	}

	return nil
}
