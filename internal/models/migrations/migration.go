package migrations

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []any {
	return []any{
		// models.IntegrationBills{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{
		// {
		// 	Model: models.OrgUserManagement{},
		// 	TableName: "org_user_managements",
		// 	Column: "is_deactivated",
		// 	Type: "bool",
		// },
	}
}
