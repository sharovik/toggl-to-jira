package entities

import "github.com/sharovik/orm/dto"

// New generates the model object for selected model interface
func New(table string, fields []interface{}, primary dto.ModelField, modelStruct dto.ModelInterface) dto.ModelInterface {
	model := modelStruct
	model.SetTableName(table)
	model.SetPrimaryKey(primary)
	for _, field := range fields {
		switch f := field.(type) {
		case dto.ModelField:
			model.AddModelField(f)
		}
	}

	return model
}
