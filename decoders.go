package HTTPNav

import (
	"bufio"
	"encoding/json"
	"strings"
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

func (ht *HTTPRequest) GetBodyAsText() {

}

func (ht *HTTPRequest) GetBodyAsJson(v *any) {
	json.NewDecoder(ht.reader).Decode(v)
}
