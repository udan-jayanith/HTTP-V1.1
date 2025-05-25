package HTTPNav

import (
	"fmt"
	"net/http"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

type HTTPResponseLine struct {
	StatusCode        int
	StatusCodeMessage string
}

func (hr *HTTPResponseLine) SetStatusCode(statusCode int) {
	hr.StatusCode = statusCode
	hr.StatusCodeMessage = http.StatusText(statusCode)
}

func (hr *HTTPResponseLine) EncodeRequestLine() []byte {
	return []byte(fmt.Sprintf("HTTP/1.1 %v %s", hr.StatusCode, hr.StatusCodeMessage))
}

type HTTPResponse struct {
	ResponseLine HTTPResponseLine
	header       []byte
	body         []byte
}

func (hr *HTTPResponse) SetHeaderLine(field, value string) {
	if hr.header == nil {
		hr.header = []byte(fmt.Sprintf("%s: %s\n", field, value))
		return
	}

	hr.header = append(hr.header, []byte(fmt.Sprintf("%s: %s\r\n", field, value))...)
}

func (hr *HTTPResponse) Write(b []byte) (int, error) {
	if hr.body == nil {
		hr.body = b
		return len(b), nil
	}

	hr.body = append(hr.body, b...)
	return len(b), nil
}

func (hr *HTTPResponse) WriteAsJson(v any) error {
	jsonByt, err := jsoniter.Marshal(v)
	if err != nil {
		return err
	}

	hr.body = jsonByt
	hr.SetHeaderLine("Content-Type", "application/json")
	return err
}

func (hr *HTTPResponse) EncodeHTTPResponse() []byte {
	res := []byte{}
	res = append(res, hr.ResponseLine.EncodeRequestLine()...)
	if hr.body != nil {
		hr.SetHeaderLine("Content-Length", strconv.Itoa(len(hr.body)))
	}

	res = append(res, '\r', '\n')
	res = append(res, hr.header...)
	res = append(res, '\r', '\n')
	res = append(res, hr.body...)
	return res
}
