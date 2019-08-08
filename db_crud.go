// package der

package der

// implement DBCrud
// form data
type DBCrud struct {
	DBCrudCom
}

// create
func (c *DBCrud) Create(params map[string][]string) error {

	return c.DBTool.CreateData(c.Param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params map[string][]string) (ID, error) {

	return c.DBTool.CreateDataResID(c.Param.Table, params)
}

// update
func (c *DBCrud) Update(params map[string][]string) error {

	return c.DBTool.UpdateData(c.Param.Table, params)
}