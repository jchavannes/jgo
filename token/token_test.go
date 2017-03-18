package token

import (
	"testing"
)

const TEST_SECRET = "test-secret"

func TestToken(t *testing.T) {

	tkn := GetSessionToken(TEST_SECRET)

	if len(tkn) == 0 {
		t.Error("No token generated")
	}

	tkn = GetSessionToken("")

	if len(tkn) == 0 {
		t.Error("Token not generated with empty secret")
	}
}

func TestValidate(t *testing.T) {

	tkn := GetSessionToken(TEST_SECRET)

	valid := Validate(tkn, TEST_SECRET)

	if valid != true {
		t.Error("Unable to validate token")
	}

	valid = Validate(tkn + "a", TEST_SECRET)

	if valid != false {
		t.Error("Incorrect token was marked as valid")
	}

	valid = Validate(tkn, TEST_SECRET + "a")

	if valid != false {
		t.Error("Different secret was marked as valid")
	}

	valid = Validate("", TEST_SECRET)

	if valid != false {
		t.Error("No token was marked as valid")
	}

	valid = Validate(tkn, "")

	if valid != false {
		t.Error("No secret was marked as valid")
	}

}
