package message

import "net/url"

type HttpRequestMessage struct {
	Method   string              `json:"method"`
	Url      string              `json:"url"`
	JsonData string              `json:"jsonData"`
	FormData url.Values          `json:"formData"`
	Sync     bool                `json:"sync"`
	Next     *HttpRequestMessage `json:"next"`
}