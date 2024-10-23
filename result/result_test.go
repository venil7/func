package result

import (
	"errors"
	"testing"
)

func TestResultOk(t *testing.T) {
	r1 := Ok("test")
	r2 := Err[any](errors.New("Bad"))

	if data, err := r1.Tuple(); err != nil || data != "test" {
		t.Fatalf(`Ok("test") = %q, %v, want "test", nil`, data, err)
	}

	if data2, err := r2.Tuple(); err == nil || data2 != nil {
		t.Fatalf(`Err[any](errors.New("Bad")) = %q, %v, want nil, non-nil`, data2, err)
	}
}
