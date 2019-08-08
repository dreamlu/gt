// package der

package der

// implement DBCrud
// form data
type DBCrudJ struct {
	DBCrudCom
}

// create
func (c *DBCrudJ) Create(data interface{}) error {

	return c.DBTool.CreateDataJ(data)
}

// create
func (c *DBCrudJ) CreateMoreDataJ(data interface{}) error {

	return c.DBTool.CreateMoreDataJ(c.Param.Table, c.Param.Model, data)
}

// update
func (c *DBCrudJ) Update(data interface{}) error {

	return c.DBTool.UpdateDataJ(data)
}
