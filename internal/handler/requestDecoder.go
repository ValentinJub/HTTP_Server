package handler

import (
	"fmt"
	"regexp"
)

type HTTPRequestDecoder interface {
	Decode(request string) *HTTPRequest
}

type HTTPRequestDecoderImpl struct{}

type RequestLine struct {
	Method      string
	RequestURI  string
	HTTPVersion string
}

func NewRequestLine(method, requestURI, httpVersion string) RequestLine {
	return RequestLine{
		Method:      method,
		RequestURI:  requestURI,
		HTTPVersion: httpVersion,
	}
}

type HTTPRequest struct {
	Request RequestLine
	Headers map[string]string
	Body    string
}

func (h *HTTPRequest) String() string {
	return h.Request.Method + " " + h.Request.RequestURI + " " + h.Request.HTTPVersion + "\r\n" +
		"Headers: " + fmt.Sprintf("%v", h.Headers) + "\r\n" +
		"Body: " + h.Body
}

func NewRequestDecoder() *HTTPRequestDecoderImpl {
	return &HTTPRequestDecoderImpl{}
}

// https://datatracker.ietf.org/doc/html/rfc2616#section-5
// Request       = Request-Line              ; Section 5.1
//                     *(( general-header        ; Section 4.5
//                      | request-header         ; Section 5.3
//                      | entity-header ) CRLF)  ; Section 7.1
//                     CRLF
//                     [ message-body ]          ; Section 4.3

// Request-Line   = Method SP Request-URI SP HTTP-Version CRLF

const (
	Request_Line_Regex = `^([A-Z]+) ([^ ]+) (HTTP\/1\.[01])\r\n`
	Header_Regex       = `^([A-z-0-9]+): (.+)`
	Body_Regex         = `\r\n\r\n(.*)$`
	CRLF               = "\r\n"
)

func (h *HTTPRequestDecoderImpl) Decode(request string) *HTTPRequest {
	requestCopy := request

	requestRegExp := regexp.MustCompile(Request_Line_Regex)
	matches := requestRegExp.FindStringSubmatch(requestCopy)
	if len(matches) < 4 {
		println("Couldn't decode HTTP Request: invalid request line format:", request)
		return nil
	}
	requestLine := NewRequestLine(matches[1], matches[2], matches[3])

	requestCopy = requestRegExp.ReplaceAllString(requestCopy, "")

	headerRegExp := regexp.MustCompile(Header_Regex)
	headers := make(map[string]string)
	for _, line := range regexp.MustCompile(CRLF).Split(requestCopy, -1) {
		if line == "" {
			break
		}
		headerMatches := headerRegExp.FindStringSubmatch(line)
		if len(headerMatches) < 3 {
			println("Couldn't decode HTTP Request: invalid header format:", line)
			return nil
		}
		headers[headerMatches[1]] = headerMatches[2]
	}

	bodyRegExp := regexp.MustCompile(Body_Regex)
	bodyMatches := bodyRegExp.FindStringSubmatch(requestCopy)
	if len(bodyMatches) < 2 {
		println("Couldn't decode HTTP Request: invalid body format:", requestCopy)
		return nil
	}
	body := bodyMatches[1]

	return &HTTPRequest{
		Request: requestLine,
		Headers: headers,
		Body:    body,
	}
}
