package qrcodegenerator

import (
	"fmt"

	"github.com/skip2/go-qrcode"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

var _ port.QRCodeGenerator = (*QRCodeGenerator)(nil)

const (
	pngSize = 256
)

type QRCodeGenerator struct {
}

func New() *QRCodeGenerator {
	return &QRCodeGenerator{}
}

func (q *QRCodeGenerator) GeneratePNG(content string) ([]byte, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, pngSize)
	if err != nil {
		return nil, fmt.Errorf("qrcode encode: %w", err)
	}

	return png, nil
}
