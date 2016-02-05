package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type (
	_upload struct {
		form string
	}

	Upload _upload
)

const (
	UPLOAD_HEADER_MD5 = "Content-MD5"
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

func (u *Upload) handleGet(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.New("upload").Parse(u.form)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(w, map[string]string{
		"action": Url(u.Path()).String(),
		"token":  u.createToken(),
	}); err != nil {
		panic(err)
	}
}

func (u *Upload) handleUpload(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(64 << 20)
	file, handler, err := req.FormFile("uploadfile")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	if !path.IsAbs(handler.Filename) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "filename: %s is not an absolute path", handler.Filename)
		return
	}

	fmt.Fprintf(w, "%v", handler.Header)

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	io.Copy(f, file)

	checksum := req.Header.Get(UPLOAD_HEADER_MD5)
	if checksum != "" {
		deleteFile := func() {
			os.Remove(handler.Filename)
		}

		code, stdout, stderr := PluginShell.shell("/usr/bin/md5sum %s", handler.Filename)
		if code != 0 {
			defer deleteFile()
			panic(fmt.Errorf("cannot calculate the MD5, %s, %s", stderr, stdout))
		}

		if !strings.Contains(stdout, checksum) {
			defer deleteFile()
			panic(fmt.Errorf("MD5 not matching. Expected: %s, but %s", checksum, stdout))
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (u *Upload) Handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		u.handleGet(w, req)
	} else {
		u.handleUpload(w, req)
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
