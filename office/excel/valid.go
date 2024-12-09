package excel

import (
	"errors"
	"github.com/dreamlu/gt/src/util"
)

func (f *Excel[T]) ValidTitle(dst *Excel[T]) (e *Excel[T], err error) {
	if len(f.Titles) != len(f.Titles) || !util.Equal(dst.Titles, f.Titles) {
		return f, errors.New("the title is different")
	}
	return
}
