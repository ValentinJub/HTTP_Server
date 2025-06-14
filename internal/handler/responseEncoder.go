package handler

import "strconv"

type Encoder interface {
	Encode(response *HTTPResponse) string
}

type HTTPResponseEncoder struct{}

func NewResponseEncoder() *HTTPResponseEncoder {
	return &HTTPResponseEncoder{}
}

type StatusLine struct {
	HTTPVersion string
	StatusCode  int
	StatusText  string
}

func NewStatusLine(httpVersion string, statusCode int, statusText string) StatusLine {
	return StatusLine{
		HTTPVersion: httpVersion,
		StatusCode:  statusCode,
		StatusText:  statusText,
	}
}

type HTTPResponse struct {
	StatusLine StatusLine
	Headers    map[string]string
	Body       string
}

func NewHTTPResponse(statusLine StatusLine, headers map[string]string, body string) *HTTPResponse {
	return &HTTPResponse{
		StatusLine: statusLine,
		Headers:    headers,
		Body:       body,
	}
}

func (e *HTTPResponseEncoder) Encode(response *HTTPResponse) string {
	responseStr := response.StatusLine.HTTPVersion + " " +
		strconv.Itoa(response.StatusLine.StatusCode) + " " +
		response.StatusLine.StatusText + CRLF

	for key, value := range response.Headers {
		responseStr += key + ": " + value + CRLF
	}

	responseStr += CRLF + response.Body

	return responseStr
}
