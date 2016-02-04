package client

import (
	"fmt"
	"net/url"
)

type ApiUrl struct {
	path  string
	query map[string]string
}

func Url(path string) *ApiUrl {
	return &ApiUrl{
		path:  path,
		query: make(map[string]string),
	}
}

func (u *ApiUrl) Query(key, value string) *ApiUrl {
	u.query[key] = value
	return u
}

func (u *ApiUrl) String() string {
	l, err := url.Parse(fmt.Sprintf("%s:///v1%s", SCHEME, u.path))
	if err != nil {
		panic(err)
	}

	if len(u.query) > 0 {
		params := url.Values{}
		for k, v := range u.query {
			params.Add(k, v)
		}

		l.RawQuery = params.Encode()
	}

	return l.String()
}
