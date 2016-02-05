package server

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"org.zstack/lib"
	"os"
	"path"
	"strconv"
	//"strings"
	//"text/template"
	"time"
)

type (
	_upload struct {
		form string
	}

	Upload _upload
)

const (
	UPLOAD_HEADER_MD5            = "Content-MD5"
	UPLOAD_HEADER_LOCATION       = "Location"
	UPLOAD_HEADER_DIGEST         = "Content-Digest"
	UPLOAD_HEADER_URL            = "URL"
	UPLOAD_HEADER_TOKEN          = "Token"
	UPLOAD_HEADER_RANGE          = "Content-Range"
	UPLOAD_HEADER_CONTENT_LENGTH = "Content-Length"

	UPLOAD_START      = iota
	UPLOAD_MONOLITHIC = iota
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
	return []string{"GET", "POST", "PUT", "DELETE", "HEAD"}
}

func (u *Upload) Path() []string {
	return []string{
		"/upload",
		"/upload/{token}",
	}
}

func (u *Upload) toAction(req *http.Request) int {
	vars := mux.Vars(req)
	_, hasToken := vars["token"]
	digest := req.URL.Query().Get("digest")
	contentRange := req.Header.Get(UPLOAD_HEADER_RANGE)

	if req.Method == "POST" && !hasToken {
		return UPLOAD_START
	} else if req.Method == "POST" && hasToken && digest != "" && contentRange == "" {
		return UPLOAD_MONOLITHIC
	} else {
		panic("unknown action")
	}
}

func (u *Upload) makeTempFilePath(location, token string) string {
	return path.Join(path.Dir(location), fmt.Sprintf("%s.upload", token))
}

func (u *Upload) start(w http.ResponseWriter, req *http.Request) {
	location := req.Header.Get(UPLOAD_HEADER_LOCATION)
	if location == "" {
		panic(fmt.Errorf("'%s' must be set in the header for starting a upload", UPLOAD_HEADER_LOCATION))
	} else if !path.IsAbs(location) {
		panic(fmt.Errorf("%s[%s] must be an absolute path", UPLOAD_HEADER_LOCATION, location))
	}

	digest := req.Header.Get(UPLOAD_HEADER_DIGEST)
	if digest == "" {
		panic(fmt.Errorf("'%s' must be set in the header for starting a upload", UPLOAD_HEADER_DIGEST))
	}

	if lib.IsFile(location) {
		if sum := lib.Md5(location); sum == digest {
			// for existing and same file, don't upload it
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	if err := os.MkdirAll(path.Dir(location), os.ModeDir); err != nil {
		panic(fmt.Errorf("cannot create directory for the file[%s], %v", location, err))
	}

	token := lib.RandomString()
	tmpFile := u.makeTempFilePath(location, token)
	f, err := os.Create(tmpFile)
	if err != nil {
		panic(fmt.Errorf("cannot create the file[%s], %v", tmpFile, err))
	}
	f.Close()

	w.Header().Add(UPLOAD_HEADER_URL, Url(fmt.Sprintf("/upload/%s", token)).String())
	w.Header().Add(UPLOAD_HEADER_TOKEN, token)
	w.WriteHeader(http.StatusAccepted)
}

func (u *Upload) monolithicUpload(w http.ResponseWriter, req *http.Request) {
	length := req.Header.Get(UPLOAD_HEADER_CONTENT_LENGTH)
	if length == "" {
		panic(fmt.Errorf("'%s' must be set in the header for a monolithic upload", UPLOAD_HEADER_CONTENT_LENGTH))
	}

	nlen, err := strconv.ParseInt(length, 10, 64)

	location := req.Header.Get(UPLOAD_HEADER_LOCATION)
	if location == "" {
		panic(fmt.Errorf("'%s' must be set in the header for a monolithic upload", UPLOAD_HEADER_LOCATION))
	}

	vars := mux.Vars(req)
	token, _ := vars["token"]
	digest := req.URL.Query().Get("digest")

	tmpFile := u.makeTempFilePath(location, token)
	if !lib.IsFile(tmpFile) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot find the temporary file at %s, have you started a upload yet?", tmpFile)
		return
	}

	f, err := os.OpenFile(tmpFile, os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Errorf("cannot open the file[%s], %v", tmpFile, err))
	}

	deleteFile := func() {
		os.Remove(tmpFile)
	}

	n, err := io.Copy(f, req.Body)
	if err != nil {
		defer deleteFile()
		panic(fmt.Errorf("failed to write to %s, %v", tmpFile, err))
	}

	if n != nlen {
		defer deleteFile()
		panic(fmt.Errorf("the file size[%d] is not matching '%s'[%d]", n, UPLOAD_HEADER_CONTENT_LENGTH, nlen))
	}

	if err := os.Rename(tmpFile, location); err != nil {
		defer deleteFile()
		panic(fmt.Errorf("failed rename the file[%s] to %s, %v", tmpFile, location, err))
	}

	sum := lib.Md5(location)
	if sum != digest {
		defer deleteFile()
		panic(fmt.Errorf("corrupted file; expected MD5[%s] but got %s", digest, sum))
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *Upload) Handler(w http.ResponseWriter, req *http.Request) {
	action := u.toAction(req)

	if action == UPLOAD_START {
		u.start(w, req)
	} else if action == UPLOAD_MONOLITHIC {
		u.monolithicUpload(w, req)
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
