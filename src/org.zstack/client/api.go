package client

import (
	"fmt"
)

func ApiURL(url string, params ...interface{}) string {
	url = fmt.Sprintf(url, params...)
	return fmt.Sprintf("/v1%s", url)
}
