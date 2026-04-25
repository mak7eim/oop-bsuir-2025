package auth

import (
	"golang.org/x/crypto/bcrypt"

	apperrors "lab7-oop/shared/errors"
)

type BcryptHasher struct{}

func (BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (BcryptHasher) Compare(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return apperrors.ErrInvalidCredentials
	}
	return nil
}
