package HTTPNav

import (
	"bufio"
	"errors"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func decodeRequestLine(reader *bufio.Reader) (RequestLine, error) {
	requestLine := RequestLine{}
	line, err := reader.ReadString('\n')
	if err != nil {
		return requestLine, err
	}

	var lastSanedIndex int
	//this decode the HTTP method
	httpMethod := ""
	for i, char := range line {
		lastSanedIndex = i
		if char == ' ' {
			break
		}

		httpMethod += string(char)
	}
	requestLine.Method = HTTPRequestMethod(strings.TrimSpace(httpMethod))

	//this decode the request target
	lastSanedIndex++
	url := RequestTarget{}
	lastSanedIndex = url.parse(line, lastSanedIndex)
	requestLine.Target = url

	//this decode the protocol
	requestLine.Protocol = strings.TrimSpace(line[lastSanedIndex:])

	return requestLine, err
}

func decodeHeader(reader *bufio.Reader) (map[string]string, error) {
	header := make(map[string]string, 9)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return header, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return header, err
		}

		for i, char := range line {
			if char == ':' {
				header[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i+1:])
				break
			}
		}
	}
}

var (
	ContentLengthHeaderNotFound error = errors.New("Expected Content-Length header but not found.")
	ContentLengthIsEmpty        error = errors.New("Content-Length header is empty.")
)

//This return Body as a slice of bytes.
//If Content-Length header is not included in the header GetBodyAsBytes() returns ContentLengthHeaderNotFound.
func (ht *HTTPRequest) GetBodyAsBytes() ([]byte, error) {
	ContentLength, ok := ht.GetHeader("Content-Length")
	ContentLength = strings.TrimSpace(ContentLength)
	if !ok {
		return []byte{}, ContentLengthHeaderNotFound
	} else if ContentLength == "" {
		return []byte{}, ContentLengthIsEmpty
	}

	contentLengthNo, err := strconv.Atoi(ContentLength)
	if err != nil {
		return []byte{}, err
	}

	bodyByt, err := ht.reader.Peek(contentLengthNo)
	return bodyByt, err
}

//GetBodyAsString return body as string. 
//It uses GetBodyAsBytes in underlying layer.
func (ht *HTTPRequest) GetBodyAsString() (string, error) {
	bodyByt, err := ht.GetBodyAsBytes()
	return string(bodyByt), err
}

var (
	ContentTypeHeaderNotFound error =errors.New("Content-Type header not found.")
	InvalidContentType error = errors.New("Invalid content-type value")
)

//GetBodyAsJson parses the JSON-encoded data and stores the result in the value pointed to by v.
//If Content-Type is not in the header GetBodyAsJson returns ContentTypeHeaderNotFound.
//Else if Content-Type != "application/json" GetBodyAsJson returns InvalidContentType.
//It uses GetBodyAsBytes in underlying layer.
func (ht *HTTPRequest) GetBodyAsJson(v any) error {
	contentType, ok := ht.GetHeader("Content-Type")
	if !ok {
		return ContentTypeHeaderNotFound
	}else if contentType != "application/json" {
		return InvalidContentType
	}

	bodyBty, err := ht.GetBodyAsBytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bodyBty, v)
	return err
}
