package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("mypassword123", 4) // low cost for test speed
	if err != nil {
		t.Fatal(err)
	}

	if err := CheckPassword(hash, "mypassword123"); err != nil {
		t.Errorf("correct password failed: %v", err)
	}

	if err := CheckPassword(hash, "wrongpassword"); err == nil {
		t.Error("wrong password should fail")
	}
}
