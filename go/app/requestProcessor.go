package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
)

func buildResponse(contentType string, statusCodeResponse int, acceptEncoding string, body string) Response {
	headers := map[string]string{}
	headers[HeaderContentLengthKey] = strconv.Itoa(len(body))
	headers[HeaderContentTypeKey] = contentType

	if len(acceptEncoding) > 0 {
		headers[HeaderContentEncodingKey] = acceptEncoding
	}

	headers[HeaderConnectionKey] = HeaderConnectionClose

	responseStatusText := ResponseStatusTextOk
	switch statusCodeResponse {
	case ResponseStatusCodeCreated:
		responseStatusText = ResponseStatusTextCreated
		break
	case ResponseStatusCodeNotFound:
		responseStatusText = ResponseStatusTextNotFound
		break
	}

	response := Response{
		protocol:   Protocol,
		statusCode: statusCodeResponse,
		statusText: responseStatusText,
		headers:    headers,
		body:       body,
	}

	return response
}

func writeResponseObject(conn net.Conn, response Response) {
	fmt.Fprintf(conn, "%s %d %s%s", response.protocol, response.statusCode, response.statusText, CRLF)
	for key, value := range response.headers {
		fmt.Fprintf(conn, "%s: %s%s", key, value, CRLF)
	}
	fmt.Fprintf(conn, "%s%s", CRLF, response.body)
}

func getAcceptEncodingHeader(headersMap map[string]string) string {
	responseAcceptEncoding := ""

	receivedAcceptEncoding := headersMap[HeaderAcceptEncodingKey]
	var acceptEncodingValues []string

	if strings.Contains(receivedAcceptEncoding, ",") {
		acceptEncodingValues = strings.Split(receivedAcceptEncoding, ",")
	}

	if len(acceptEncodingValues) > 0 {
		for _, element := range acceptEncodingValues {
			trimElement := strings.TrimSpace(element)
			if slices.Contains(SupportedEncodings, trimElement) {
				if len(responseAcceptEncoding) == 0 {
					responseAcceptEncoding = trimElement
				} else {
					responseAcceptEncoding = responseAcceptEncoding + ", " + trimElement
				}
			}
		}
	} else {
		if slices.Contains(SupportedEncodings, receivedAcceptEncoding) {
			responseAcceptEncoding = receivedAcceptEncoding
		}
	}

	return responseAcceptEncoding
}

func processGetRequest(conn net.Conn, request Request) {
	if len(request.path) == 1 && strings.EqualFold(request.path, "/") {
		response := buildResponse(HeaderContentTypeText, ResponseStatusCodeOk, "", "")
		writeResponseObject(conn, response)
	} else if strings.Index(request.path, "/echo") == 0 {
		splitPath := strings.Split(request.path, "/")
		receivedStr := splitPath[len(splitPath)-1]

		responseAcceptEncoding := getAcceptEncodingHeader(request.headers)

		responseBody := receivedStr
		if strings.Contains(responseAcceptEncoding, "gzip") {
			responseBody, _ = writeGzipBuffer(receivedStr)
		}

		response := buildResponse(HeaderContentTypeText, ResponseStatusCodeOk, responseAcceptEncoding, responseBody)
		writeResponseObject(conn, response)
	} else if strings.Index(request.path, "/user-agent") == 0 {
		userAgent := request.headers[HeaderUserAgentKey]
		response := buildResponse(HeaderContentTypeText, ResponseStatusCodeOk, "", userAgent)
		writeResponseObject(conn, response)
	} else if strings.Index(request.path, "/files") == 0 {
		tmpDirectory := ""
		if len(os.Args) > 1 {
			tmpDirectory = os.Args[2]
		}
		splitPath := strings.Split(request.path, "/")
		fileName := splitPath[len(splitPath)-1]
		fileContent := readFile(tmpDirectory, fileName)

		headers := map[string]string{}
		var response Response

		if fileContent != nil {
			headers[HeaderContentTypeKey] = HeaderContentTypeOctetStream
			headers[HeaderContentLengthKey] = strconv.Itoa(len(fileContent))
			headers[HeaderConnectionKey] = HeaderConnectionClose
			response = Response{
				protocol:   Protocol,
				statusCode: ResponseStatusCodeOk,
				statusText: ResponseStatusTextOk,
				headers:    headers,
				body:       string(fileContent),
			}
			response = buildResponse(HeaderContentTypeOctetStream, ResponseStatusCodeOk, "", string(fileContent))

		} else {
			response = buildResponse(HeaderContentTypeText, ResponseStatusCodeNotFound, "", "")
		}
		writeResponseObject(conn, response)
	} else {
		response := buildResponse(HeaderContentTypeText, ResponseStatusCodeNotFound, "", "")
		writeResponseObject(conn, response)
	}
}

func processPostRequest(conn net.Conn, request Request, reader *bufio.Reader) {
	if strings.Index(request.path, "/files") == 0 {

		splitPath := strings.Split(request.path, "/")
		fileName := splitPath[len(splitPath)-1]

		contentType := request.headers[HeaderContentTypeKey]
		contentLengthString := strings.Replace(request.headers[HeaderContentLengthKey], CRLF, "", 1)
		fmt.Println("ContentType:", contentType)

		var response Response

		if len(contentLengthString) > 0 {
			length, err := strconv.Atoi(contentLengthString)
			if err != nil {
				response = buildResponse(HeaderContentTypeText, ResponseStatusCodeNotFound, "", "")
			}

			writeFile(length, reader, fileName)
			response = buildResponse(HeaderContentTypeText, ResponseStatusCodeCreated, "", "")
		}

		writeResponseObject(conn, response)
	}
}

func processRequest(conn net.Conn, request Request, reader *bufio.Reader) {
	if request.method == "GET" {
		processGetRequest(conn, request)
	} else if request.method == "POST" {
		processPostRequest(conn, request, reader)
	}

}

func handleConnection(conn net.Conn) {
	var request Request

	// Enable persistent connections by using the same connection to process many requests
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading request:", err)
			break
		}

		contents := strings.Split(message, " ")
		method := contents[0]
		path := contents[1]
		protocol := strings.Replace(contents[2], CRLF, "", 1)

		headersMap := map[string]string{}
		for {
			headerMsg, _ := reader.ReadString('\n')
			if headerMsg == "" || headerMsg == CRLF {
				break
			}

			headerMsgSplit := strings.Split(headerMsg, ":")
			headerKey := headerMsgSplit[0]
			headerValue := headerMsgSplit[1]
			headersMap[headerKey] = strings.TrimSpace(headerValue)
		}

		request = Request{
			method:   method,
			path:     path,
			protocol: protocol,
			headers:  headersMap,
			body:     "",
		}
		fmt.Println("Received Request: ", request)

		processRequest(conn, request, reader)

		if request.headers[HeaderConnectionKey] == HeaderConnectionClose {
			conn.Close()
		}
	}
}
