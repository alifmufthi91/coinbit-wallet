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
		// var result T
		*dest = *new(T)
		return nil
	}
	var ok bool
	*dest, ok = val.(T)
	if !ok {
		return errors.New("failed to cast type")
	}
	return nil
}
