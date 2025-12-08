package convertor

import "encoding/base64"

func encodeBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
