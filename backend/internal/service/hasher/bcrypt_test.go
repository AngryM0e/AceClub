package hasher

import (
	"testing"
)

func TestBcrypt_Success(t *testing.T) {
	hasher, err := NewBcryptHasher(8)
	if err != nil {
		t.Errorf("wrong cost: %v", err)
	}

	pass, err := hasher.Hash("abcde")
	if err != nil {
		t.Error(err)
	}

	err = hasher.Compare(pass, "abcde")
	if err != nil {
		t.Error(err)
	}
}

func TestBcrype_InvalidCost(t *testing.T) {
	_, err := NewBcryptHasher(3)
	if err == nil {
		t.Error("expected error for cost=3, got nil")
	}

	_, err = NewBcryptHasher(32)
	if err == nil {
		t.Error("expected error for cost=32, got nil")
	}
}

func TestBcrypt_WrongPassword(t *testing.T) {
	hasher, _ := NewBcryptHasher(8)
	pass, _ := hasher.Hash("abcde")
	err := hasher.Compare(pass, "edcba")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}
