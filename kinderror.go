package cliconfig

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type KindError struct {
	Unsupported reflect.Kind
	MustBe      []reflect.Kind
	MetaKind    *reflect.Kind
}

func (k KindError) Is(target error) bool {
	return true
}

func newKindError(unsupported reflect.Kind, mustBe []reflect.Kind, msg string) error {
	return errors.Wrap(KindError{Unsupported: unsupported, MustBe: mustBe}, msg)
}

func (t KindError) Error() string {
	if t.MetaKind != nil {
		return fmt.Sprintf("%v kind of %v is not supported; expecting %v", t.MetaKind, t.Unsupported, t.MustBe)
	}
	return fmt.Sprintf("kind %v is not supported; expecting %v", t.Unsupported, t.MustBe)
}
