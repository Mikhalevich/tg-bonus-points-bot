package request

import (
	"errors"
)

type OrderStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (os OrderStatus) Validate() error {
	if os.ID == "" {
		return errors.New("empty id")
	}

	if os.Status == "" {
		return errors.New("empty status")
	}

	return nil
}
