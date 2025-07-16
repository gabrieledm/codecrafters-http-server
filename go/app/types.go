package main

type Request struct {
	method   string
	path     string
	protocol string
	headers  map[string]string
	body     string
}

type Response struct {
	protocol   string
	statusCode int
	statusText string
	headers    map[string]string
	body       string
}
