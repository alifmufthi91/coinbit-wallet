package util

import (
	"errors"

	"github.com/lovoo/goka"
)

func GetView[T any](view *goka.View, stream string, dest *T) error {
	val, err := view.Get(stream)
	if err != nil {
		return err
	} else if val == nil {
		return errors.New("view is not found")
	}
	var ok bool
	*dest, ok = val.(T)
	if !ok {
		return errors.New("failed to cast type")
	}
	return nil
}
