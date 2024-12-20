package qrcodegenerator

import (
	"fmt"

	"github.com/skip2/go-qrcode"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

var _ port.QRCodeGenerator = (*QRCodeGenerator)(nil)

type QRCodeGenerator struct {
}

func New() *QRCodeGenerator {
	return &QRCodeGenerator{}
}

func (q *QRCodeGenerator) GeneratePNG(content string) ([]byte, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("qrcode encode: %w", err)
	}

	return png, nil
}
