package tokens

import (
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
)

func New() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}
