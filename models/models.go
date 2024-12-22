package models

import (
	"github.com/unvs/models/accounts"
)

type Model struct {
	Account accounts.Accounts `table:"accounts`
}
