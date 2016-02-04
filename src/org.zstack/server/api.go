package server

import (
	"fmt"
)

func ApiURL(url string) string {
	return fmt.Sprintf("/v1%s", url)
}
