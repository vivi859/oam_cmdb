package test

import (
	"OAM/models"
	"testing"
)

func TestGetAccount(t *testing.T) {
	acct := models.GetAccountById(1)
	t.Log(acct)
}
