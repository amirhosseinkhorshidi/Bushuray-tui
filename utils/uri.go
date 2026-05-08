package utils

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// ParseURIInfo extracts address, port, and transport type from a proxy URI.
func ParseURIInfo(uri string) (address, port, transport string) {
	if uri == "" {
		return
	}
	schemeEnd := strings.Index(uri, "://")
	if schemeEnd == -1 {
		return
	}
	scheme := uri[:schemeEnd]
	rest := uri[schemeEnd+3:]

	if scheme == "vmess" {
		if idx := strings.Index(rest, "#"); idx != -1 {
			rest = rest[:idx]
		}
		var decoded []byte
		for _, enc := range []*base64.Encoding{
			base64.StdEncoding,
			base64.RawStdEncoding,
			base64.URLEncoding,
			base64.RawURLEncoding,
		} {
			if d, err := enc.DecodeString(rest); err == nil {
				decoded = d
				break
			}
		}
		if decoded == nil {
			return
		}
		var cfg struct {
			Add  string      `json:"add"`
			Port interface{} `json:"port"`
			Net  string      `json:"net"`
		}
		if json.Unmarshal(decoded, &cfg) != nil {
			return
		}
		address = cfg.Add
		transport = cfg.Net
		switch p := cfg.Port.(type) {
		case float64:
			port = strconv.Itoa(int(p))
		case string:
			port = p
		}
		return
	}

	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	address = u.Hostname()
	port = u.Port()
	transport = u.Query().Get("type")
	if transport == "" {
		transport = "tcp"
	}
	return
}
