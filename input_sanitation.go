package main

import (
	"bytes"
	"strings"
)

func CleanBodyBytes(value []byte) []byte {
	stripped := bytes.ReplaceAll(value, []byte("\r"), []byte{})
	stripped = bytes.ReplaceAll(stripped, []byte("\n"), []byte{})
	stripped = bytes.ReplaceAll(stripped, []byte("\t"), []byte{})
	stripped = bytes.ReplaceAll(stripped, []byte(" "), []byte{})
	return bytes.ToLower(stripped)
}

func CleanBodyString(value string) string {
	return string(CleanBodyBytes([]byte(value)))
}

func CleanHeaders(headers map[string][]string) map[string][]string {
	var ret = make(map[string][]string)
	for key, value := range headers {
		ret[cleanHeaderKey(key)] = cleanHeaderValues(value)
	}
	return ret
}

func cleanHeaderValues(values []string) []string {
	cleaned := make([]string, len(values))
	for i, item := range values {
		cleaned[i] = cleanHeaderValue(item)
	}
	return cleaned
}

func cleanHeaderValue(header string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToLower(header),
				" ", ""),
			",", ""),
		";", "")
}

func cleanHeaderKey(header string) string {
	return strings.Trim(strings.ToLower(header), " ")
}

func CleanMethod(method string) string {
	return strings.Trim(strings.ToLower(method), " ")
}

func CleanUrl(url string) string {
	return strings.Trim(strings.ToLower(url), " /")
}
