package password

import "golang.org/x/crypto/bcrypt"

func GenerateHash(password string, result *string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	*result = string(hashedPassword)

	return nil
}
