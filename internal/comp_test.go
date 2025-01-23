package internal_test

import (
	"testing"

	"github.com/hizla/hizla/internal"
)

func TestCheck(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if _, ok := internal.Check(""); ok {
			t.FailNow()
		}
	})
	t.Run("poison", func(t *testing.T) {
		if _, ok := internal.Check("INVALIDINVALIDINVALIDINVALIDINVALID"); ok {
			t.FailNow()
		}
	})
	t.Run("valid", func(t *testing.T) {
		v := "\n"
		if got, ok := internal.Check(v); !ok || v != got {
			t.FailNow()
		}
	})
}
