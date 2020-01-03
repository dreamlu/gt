package te

// customize cn text error
type TextError struct {
	Msg string
}

func (s *TextError) Error() string {
	return s.Msg
}

var TextErr *TextError
