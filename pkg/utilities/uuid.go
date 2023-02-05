package utilities

import "github.com/google/uuid"

// NewUUIDString returns a new UUID string value
func NewUUIDString() (string, error) {
	id, e := uuid.NewUUID()
	if e != nil {
		return "", e
	}
	return id.String(), nil
}
