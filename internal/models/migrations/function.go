package migrations

type AlterColumn struct {
	Model     any
	TableName string
	Column    string
	Type      string
}

// func (a *AlterColumn) UpdateColumnType(db *gorm.DB) error {
// 	return db.Migrator().AlterColumn(a.Model, a.TableName, a.Column, a.Type)
// }
