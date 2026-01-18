package utils

import (
	"strconv"
	"strings"
)

func ConvertUint16(str string) uint16 {
	u, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(u)
}

func ConvertUint8(str string) uint8 {
	u, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(u)
}

func ExtractBearerToken(authHeader string) string {
	const prefix = "Bearer "
	tokenIndex := strings.Index(authHeader, prefix)
	if tokenIndex == -1 || tokenIndex != 0 {
		return ""
	}
	return authHeader[tokenIndex+len(prefix):]
}

func FormatEndpoint(endpoint string) string {
	endpoint = strings.ReplaceAll(endpoint, " ", "")
	endpoint = strings.ReplaceAll(endpoint, "/", "-")
	endpoint = strings.ReplaceAll(endpoint, "?", "")
	return endpoint
}
