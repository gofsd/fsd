package account

import (
	"errors"
)

var AccountNotFound = errors.New("Account Not Found")
var AccountExist = errors.New("Account Already Exists")
