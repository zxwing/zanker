package server

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type (
	_upload struct {
		form string
	}

	Upload _upload
)

var (
	PluginUpload = &Upload{
		form: `<html>
<head>
<title>Upload file</title>
</head>
	
<body>
	<form enctype="multipart/form-data" action="{{.action}}" method="post">
		<input type="file" name="uploadfile" />
		<input type="hidden" name="token" value="{{.token}}"/>
		<input type="submit" value="upload" />
	</form>
</body>
</html>`,
	}
)

func (u *Upload) Methods() []string {
	return []string{"GET", "POST"}
}

func (u *Upload) Path() string {
	return "/upload"
}

func (u *Upload) Handler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.New("upload").Parse(u.form)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{
		"action": Url(u.Path()).String(),
		"token":  u.createToken(),
	})

	if err != nil {
		panic(err)
	}
}

func (u *Upload) createToken() string {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}

func init() {
	RegisterApiRoute(PluginUpload)
}
