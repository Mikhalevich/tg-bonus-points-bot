package request

import (
	"errors"
)

type OrderStatus struct {
	Status string `json:"status"`
}

func (os OrderStatus) Validate() error {
	if os.Status == "" {
		return errors.New("empty status")
	}

	return nil
}
