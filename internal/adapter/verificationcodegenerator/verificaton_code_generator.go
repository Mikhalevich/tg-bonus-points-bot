package verificationcodegenerator

import (
	"fmt"
	"math/rand"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

const (
	maxCodeValue = 1000
)

var _ port.VerificationCodeGenerator = (*VerificationCodeGenerator)(nil)

type VerificationCodeGenerator struct {
}

func New() *VerificationCodeGenerator {
	return &VerificationCodeGenerator{}
}

func (g *VerificationCodeGenerator) Generate() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(maxCodeValue))
}
