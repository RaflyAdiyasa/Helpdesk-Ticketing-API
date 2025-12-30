package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GeneratePrefixedUUID(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, uuid.New().String())
}

func GenerateUserID(role string) string {
	return GeneratePrefixedUUID(role)
}

func GenerateTicketID() string {
	return GeneratePrefixedUUID("tick")
}
