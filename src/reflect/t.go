package reflect

// SetT data field value
// type and field must same
func SetT(data any, field string, value string) {
	TrueValueOf(data).FieldByName(field).Set(TrueValueOf(value))
}
