package jsonb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB string

func NewFromMarshaler(s any) (JSONB, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	return JSONB(b), nil
}

func (j JSONB) Value() (driver.Value, error) {
	if j == "" {
		return "{}", nil
	}

	return j, nil
}

func ConvertTo[T any](j JSONB, v *T) error {
	if err := json.Unmarshal([]byte(j), v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}
