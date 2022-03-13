package httpr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	Method string
	Path   string
	Params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request: r,
		Writer:  w,
		Method:  r.Method,
		Path:    r.URL.Path,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

// ParseJSON parse body data as json format.
func (c *Context) ParseJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(obj)
}

// ParseJSONObj parse body data as json format.
func (c *Context) ParseJSONObj(obj interface{}) (interface{}, error) {
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (c *Context) SetHeader(k, v string) {
	c.Writer.Header().Set(k, v)
}

func (c *Context) SetStatus(code int) {
	c.Writer.WriteHeader(code)
}

// ---------- Response return -------------------------------------------------------

// JSON return json format data, use application/json as content type.
func (c *Context) JSON(obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(http.StatusOK)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		c.error(http.StatusInternalServerError, err)
	}
}

// String return string, use text/plain as content type.
func (c *Context) String(format string, values ...interface{}) {
	c.StringWithHttpStatus(http.StatusOK, format, values...)
}

func (c *Context) StringWithHttpStatus(status int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(status)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Data(data []byte) {
	c.SetStatus(http.StatusOK)
	_, _ = c.Writer.Write(data)
}

// HTML return html string, use text/html as content type.
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	_, _ = c.Writer.Write([]byte(html))
}

func (c *Context) error(status int, err error) {
	c.SetStatus(status)
	_, _ = c.Writer.Write([]byte(err.Error()))
}

// Result common result
type Result struct {
	RequestId string `json:"requestId"`

	code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (c *Context) Result(message string, data interface{}, err error) {
	result := &Result{}
	if err != nil {
		result.code = -1
		result.Msg = err.Error()
	} else {
		result.Msg = message
		result.Data = data
	}
	c.JSON(result)
}

func (c *Context) ResultError(err error) {
	c.ResultErrorWithCode(-1, err)
}

func (c *Context) ResultErrorWithCode(code int, err error) {
	result := &Result{
		code: -1,
		Msg:  "internal error",
	}
	if err != nil {
		result.code = code
		result.Msg = err.Error()
	}
	c.JSON(result)
}

func (c *Context) ResultMessage(message string, err error) {
	result := &Result{}
	if err != nil {
		result.code = -1
		result.Msg = err.Error()
	} else {
		result.Msg = message
	}
	c.JSON(result)
}

func (c *Context) ResultData(data interface{}, err error) {
	result := &Result{}
	if err != nil {
		result.code = -1
		result.Msg = err.Error()
	} else {
		result.Msg = "success"
		result.Data = data
	}
	c.JSON(result)
}
