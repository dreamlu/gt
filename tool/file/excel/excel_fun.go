package excel

func ExportExcel(model, data interface{}) (e *Excel, err error) {
	e = NewExcel(model)
	err = e.Export(data)
	return
}
