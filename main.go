/*
HTTPNav is a HTTP 1.1 go package.
*/
package HTTPNav

import (
	"bufio"
	"net"
	"net/url"
	"strings"
)

type HTTPCallbackHandleFunc func(string, *HTTPRequest)

type Server struct {
	handlers map[string]HTTPCallbackHandleFunc
}

type HTTPRequestMethod string

//Http methods.
var (
	/*The GET method requests a representation of the specified resource. Requests using GET should only retrieve data and should not contain a request content.*/
	Get HTTPRequestMethod = "GET"
	/*The HEAD method asks for a response identical to a GET request, but without a response body.*/
	Head HTTPRequestMethod = "HEAD"
	/*The POST method submits an entity to the specified resource, often causing a change in state or side effects on the server.*/
	Post HTTPRequestMethod = "POST"
	/*The PUT method replaces all current representations of the target resource with the request content.*/
	Put HTTPRequestMethod = "PUT"
	/*The DELETE method deletes the specified resource.*/
	Delete HTTPRequestMethod = "DELETE"
	/*The CONNECT method establishes a tunnel to the server identified by the target resource.*/
	Connect HTTPRequestMethod = "CONNECT"
	/*The OPTIONS method describes the communication options for the target resource.*/
	Options HTTPRequestMethod = "OPTIONS"
	/*The TRACE method performs a message loop-back test along the path to the target resource.*/
	Trace HTTPRequestMethod = "TRACE"
	/*The PATCH method applies partial modifications to a resource.*/
	Patch HTTPRequestMethod = "PATCH"
)

type RequestTarget struct {
	RequestTarget string
	Path          string
	QueryParams   map[string]string
}

func (target *RequestTarget) parse(line string, startingIndex int) int {
	for ; line[startingIndex] != ' '; startingIndex++ {
		target.RequestTarget += string(line[startingIndex])
	}

	target.RequestTarget, _ = url.QueryUnescape(strings.TrimSpace(target.RequestTarget))

	target.QueryParams = make(map[string]string, 1)
	questionMarkFound := false
	key := ""
	value := ""
	equalSignFound := false
	for _, char := range target.RequestTarget {
		if char == '?' {
			questionMarkFound = true
			continue
		} else if !questionMarkFound {
			target.Path += string(char)
		} else if questionMarkFound {
			if char == '=' {
				equalSignFound = true
				continue
			} else if char == '&' {
				target.QueryParams[strings.TrimSpace(key)] = strings.TrimSpace(value)
				key = ""
				value = ""
				equalSignFound = false
				continue
			} else if !equalSignFound {
				key += string(char)
			} else if equalSignFound {
				value += string(char)
			}
		}
	}
	if strings.TrimSpace(key) != "" {
		target.QueryParams[key] = value
	}
	return startingIndex
}

//GetQuery returns RequestTarget query param value for given field.
func (target *RequestTarget) GetQuery(field string) (string, bool) {
	value, ok := target.QueryParams[field]
	return value, ok
}

type RequestLine struct {
	Method   HTTPRequestMethod
	Target   RequestTarget
	Protocol string
}

type HTTPRequest struct {
	RequestLine RequestLine
	Header      map[string]string
	reader *bufio.Reader
	isGetReaderUsed bool
}

//GetHeader returns header value for a give field.
func (ht *HTTPRequest) GetHeader(field string) (string, bool) {
	value, ok := ht.Header[field]
	return value, ok
}

//GetReader() return the reader after reading the header.
func (ht *HTTPRequest) GetReader() *bufio.Reader{
	ht.isGetReaderUsed = true
	return ht.reader
}

// GetServer returns a new Server
func GetServer() *Server {
	return &Server{
		handlers: make(map[string]HTTPCallbackHandleFunc, 18),
	}
}

// HandelFunc takes HTTPRequestMethod, requestTarget and a handler. If requests httpMethod and requestTarget matches the handler handler will execute.
func (s *Server) HandleFunc(method HTTPRequestMethod, requestTarget string, handler HTTPCallbackHandleFunc) {
	s.handlers[requestTarget] = handler
}

//StartServer start the server(start listing).
func (s *Server) StartServer(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	httpRequest := HTTPRequest{}

	requestLine, err := decodeRequestLine(reader)
	if err != nil {
		return
	}
	httpRequest.RequestLine = requestLine

	header, err := decodeHeader(reader)
	if err != nil {
		return
	}
	httpRequest.Header = header
	httpRequest.reader = reader
	callback, ok := s.handlers[httpRequest.RequestLine.Target.Path]
	if ok{
		callback("", &httpRequest)
	}
}
