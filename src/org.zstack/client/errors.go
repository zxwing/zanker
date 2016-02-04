package client

import (
	"fmt"
)

func ServerError(msg string) error {
	return fmt.Errorf("Server Error: %s", msg)
}
