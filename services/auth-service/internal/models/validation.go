package models

import (
	"fmt"
)

func (req *RegisterRequest) Validate() error {
	if req.Name == "" || req.Username == "" || req.Password == "" {
		return fmt.Errorf("Username, name and password are required")
	}

	if len(req.Password) < 6 {
		return fmt.Errorf("The password must contain at least 6 characters")
	}

	return nil
}
