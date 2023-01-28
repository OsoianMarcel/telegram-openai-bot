package stats

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
)

// Generate a string hash from user id.
func HashUserId(userId int64) string {
	userIdStr := strconv.FormatInt(userId, 10)
	hashBytes := sha1.Sum([]byte(userIdStr))

	return hex.EncodeToString(hashBytes[:])
}
