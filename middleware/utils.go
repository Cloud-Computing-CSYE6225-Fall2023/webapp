package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/shivasaicharanruthala/webapp/errors"
)

func getEmailAndPassword(authHeader string) (email, pass string, err error) {
	const authLen = 2
	auth := strings.SplitN(authHeader, " ", authLen)

	if authHeader == "" {
		return "", "", errors.NewCustomError(errors.ErrMissingHeader, 401)
	}

	if len(auth) != authLen || auth[0] != "Basic" {
		return "", "", errors.NewCustomError(errors.ErrInvalidHeader, 401)
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", authLen)

	if len(pair) < authLen {
		return "", "", errors.NewCustomError(errors.ErrInvalidToken, 401)
	}

	return pair[0], pair[1], nil
}
