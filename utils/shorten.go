package utils

import (
	"encoding/base64"
	"fmt"
	"time"
)

func GetShortCode() string {
	timestamp := time.Now().UnixNano()
	timestampBytes := []byte(fmt.Sprintf("%d", timestamp))

	key := base64.URLEncoding.EncodeToString(timestampBytes)
	key = key[:len(key)-2]
	return key[16:]
}
