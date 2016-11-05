package models

import "github.com/cSploit/daemon/models/internal"

func Setup() error {
	return internal.OpenDb()
}
