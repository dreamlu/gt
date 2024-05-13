package excel

import (
	"errors"
	"github.com/dreamlu/gt/src/util"
)

func (f *Excel[T]) ValidTitle(dst *Excel[T]) (e *Excel[T], err error) {
	if len(dst.rows[0]) != len(f.rows[0]) || !util.Equal(dst.rows[0], f.rows[0]) {
		return f, errors.New("the title is different")
	}
	return
}
