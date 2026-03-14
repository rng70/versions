package vars

import (
	"errors"
	"testing"
)

// ─── ErrUnsupportedSource ─────────────────────────────────────────────────────

func TestErrUnsupportedSource_IsNotNil(t *testing.T) {
	if ErrUnsupportedSource == nil {
		t.Error("ErrUnsupportedSource should not be nil")
	}
}

func TestErrUnsupportedSource_MessageContent(t *testing.T) {
	want := "unsupported version source (URL or file)"
	if ErrUnsupportedSource.Error() != want {
		t.Errorf("got %q, want %q", ErrUnsupportedSource.Error(), want)
	}
}

func TestErrUnsupportedSource_MessageNonEmpty(t *testing.T) {
	if ErrUnsupportedSource.Error() == "" {
		t.Error("ErrUnsupportedSource message should not be empty")
	}
}

func TestErrUnsupportedSource_IsErrorType(t *testing.T) {
	var err error = ErrUnsupportedSource
	if err == nil {
		t.Error("ErrUnsupportedSource should satisfy the error interface")
	}
}

func TestErrUnsupportedSource_CanBeCompared(t *testing.T) {
	var err error = ErrUnsupportedSource
	if !errors.Is(err, ErrUnsupportedSource) {
		t.Error("errors.Is should match ErrUnsupportedSource")
	}
}

func TestErrUnsupportedSource_NotEqualToOtherError(t *testing.T) {
	other := errors.New("some other error")
	if errors.Is(other, ErrUnsupportedSource) {
		t.Error("a different error should not match ErrUnsupportedSource")
	}
}
