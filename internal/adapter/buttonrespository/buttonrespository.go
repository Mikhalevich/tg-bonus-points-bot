package buttonrespository

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

var (
	_ port.ButtonRepositoryWriter = (*ButtonRepository)(nil)
	_ port.ButtonRepositoryReader = (*ButtonRepository)(nil)
)

type ButtonRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *ButtonRepository {
	return &ButtonRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *ButtonRepository) IsNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}

func generateID() string {
	return uuid.NewString()
}

func encodeButton(b button.Button) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func decodeButton(b []byte) (*button.Button, error) {
	var btn button.Button
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&btn); err != nil {
		return nil, fmt.Errorf("gob decode: %w", err)
	}

	return &btn, nil
}

func parseButtonID(id button.ID) (string, string) {
	var (
		idStr          = id.String()
		keySplitterIdx = strings.Index(idStr, "_")
	)

	if keySplitterIdx == -1 {
		return idStr, ""
	}

	return idStr[:keySplitterIdx], idStr[keySplitterIdx+1:]
}

func makeHmapbuttonID(key string, num string) button.ID {
	return button.IDFromString(fmt.Sprintf("%s_%s", key, num))
}
