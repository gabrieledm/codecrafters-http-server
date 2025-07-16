package main

const Protocol = "HTTP/1.1"
const HeaderContentTypeKey = "Content-Type"
const HeaderContentLengthKey = "Content-Length"
const HeaderContentEncodingKey = "Content-Encoding"
const HeaderConnectionKey = "Connection"
const HeaderUserAgentKey = "User-Agent"
const HeaderAcceptEncodingKey = "Accept-Encoding"
const CRLF = "\r\n"

const HeaderContentTypeText = "text/plain"
const HeaderContentTypeOctetStream = "application/octet-stream"
const HeaderConnectionClose = "close"

const ResponseStatusCodeOk = 200
const ResponseStatusTextOk = "OK"
const ResponseStatusCodeCreated = 201
const ResponseStatusTextCreated = "Created"
const ResponseStatusCodeNotFound = 404
const ResponseStatusTextNotFound = "Not Found"

var SupportedEncodings = []string{"gzip"}
