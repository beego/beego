package context

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type BeegoOutput struct {
	Context    *Context
	Status     int
	EnableGzip bool
	res        http.ResponseWriter
}

func NewOutput(res http.ResponseWriter) *BeegoOutput {
	return &BeegoOutput{
		res: res,
	}
}

func (output *BeegoOutput) Header(key, val string) {
	output.res.Header().Set(key, val)
}

func (output *BeegoOutput) Body(content []byte) {
	output_writer := output.res.(io.Writer)
	if output.EnableGzip == true && output.Context.Input.Header("Accept-Encoding") != "" {
		splitted := strings.SplitN(output.Context.Input.Header("Accept-Encoding"), ",", -1)
		encodings := make([]string, len(splitted))

		for i, val := range splitted {
			encodings[i] = strings.TrimSpace(val)
		}
		for _, val := range encodings {
			if val == "gzip" {
				output.Header("Content-Encoding", "gzip")
				output_writer, _ = gzip.NewWriterLevel(output.res, gzip.BestSpeed)

				break
			} else if val == "deflate" {
				output.Header("Content-Encoding", "deflate")
				output_writer, _ = flate.NewWriter(output.res, flate.BestSpeed)
				break
			}
		}
	} else {
		output.Header("Content-Length", strconv.Itoa(len(content)))
	}
	output_writer.Write(content)
	switch output_writer.(type) {
	case *gzip.Writer:
		output_writer.(*gzip.Writer).Close()
	case *flate.Writer:
		output_writer.(*flate.Writer).Close()
	case io.WriteCloser:
		output_writer.(io.WriteCloser).Close()
	}
}

func (output *BeegoOutput) Cookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	if len(others) > 0 {
		switch others[0].(type) {
		case int:
			if others[0].(int) > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", others[0].(int))
			} else if others[0].(int) < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int64:
			if others[0].(int64) > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", others[0].(int64))
			} else if others[0].(int64) < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int32:
			if others[0].(int32) > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", others[0].(int32))
			} else if others[0].(int32) < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		}
	}
	if len(others) > 1 {
		fmt.Fprintf(&b, "; Path=%s", sanitizeValue(others[1].(string)))
	}
	if len(others) > 2 {
		fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(others[2].(string)))
	}
	if len(others) > 3 {
		fmt.Fprintf(&b, "; Secure")
	}
	if len(others) > 4 {
		fmt.Fprintf(&b, "; HttpOnly")
	}
	output.res.Header().Add("Set-Cookie", b.String())
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

func (output *BeegoOutput) Json(data interface{}, hasIndent bool, coding bool) error {
	output.Header("Content-Type", "application/json;charset=UTF-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		http.Error(output.res, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding {
		content = []byte(stringsToJson(string(content)))
	}
	output.Body(content)
	return nil
}

func (output *BeegoOutput) Jsonp(data interface{}, hasIndent bool) error {
	output.Header("Content-Type", "application/javascript;charset=UTF-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		http.Error(output.res, err.Error(), http.StatusInternalServerError)
		return err
	}
	callback := output.Context.Input.Query("callback")
	if callback == "" {
		return errors.New(`"callback" parameter required`)
	}
	callback_content := bytes.NewBufferString(template.JSEscapeString(callback))
	callback_content.WriteString("(")
	callback_content.Write(content)
	callback_content.WriteString(");\r\n")
	output.Body(callback_content.Bytes())
	return nil
}

func (output *BeegoOutput) Xml(data interface{}, hasIndent bool) error {
	output.Header("Content-Type", "application/xml;charset=UTF-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = xml.MarshalIndent(data, "", "  ")
	} else {
		content, err = xml.Marshal(data)
	}
	if err != nil {
		http.Error(output.res, err.Error(), http.StatusInternalServerError)
		return err
	}
	output.Body(content)
	return nil
}

func (output *BeegoOutput) Download(file string) {
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Disposition", "attachment; filename="+filepath.Base(file))
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "must-revalidate")
	output.Header("Pragma", "public")
	http.ServeFile(output.res, output.Context.Request, file)
}

func (output *BeegoOutput) ContentType(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		output.Header("Content-Type", ctype)
	}
}

func (output *BeegoOutput) SetStatus(status int) {
	output.res.WriteHeader(status)
	output.Status = status
}

func (output *BeegoOutput) IsCachable(status int) bool {
	return output.Status >= 200 && output.Status < 300 || output.Status == 304
}

func (output *BeegoOutput) IsEmpty(status int) bool {
	return output.Status == 201 || output.Status == 204 || output.Status == 304
}

func (output *BeegoOutput) IsOk(status int) bool {
	return output.Status == 200
}

func (output *BeegoOutput) IsSuccessful(status int) bool {
	return output.Status >= 200 && output.Status < 300
}

func (output *BeegoOutput) IsRedirect(status int) bool {
	return output.Status == 301 || output.Status == 302 || output.Status == 303 || output.Status == 307
}

func (output *BeegoOutput) IsForbidden(status int) bool {
	return output.Status == 403
}

func (output *BeegoOutput) IsNotFound(status int) bool {
	return output.Status == 404
}

func (output *BeegoOutput) IsClientError(status int) bool {
	return output.Status >= 400 && output.Status < 500
}

func (output *BeegoOutput) IsServerError(status int) bool {
	return output.Status >= 500 && output.Status < 600
}

func stringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}
