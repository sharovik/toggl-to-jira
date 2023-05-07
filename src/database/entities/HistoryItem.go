package entities

import "github.com/sharovik/orm/dto"

type HistoryModelStruct struct {
	dto.BaseModel
}

func NewHistoryItem() dto.ModelInterface {
	return New(
		"history",
		[]interface{}{
			dto.ModelField{
				Name:   "task_key",
				Type:   dto.VarcharColumnType,
				Length: 255,
			},
			dto.ModelField{
				Name: "duration",
				Type: dto.IntegerColumnType,
			},
			dto.ModelField{
				Name:   "added",
				Type:   dto.VarcharColumnType,
				Length: 255,
			},
		},
		dto.ModelField{
			Name:          "id",
			Type:          dto.IntegerColumnType,
			AutoIncrement: true,
			IsPrimaryKey:  true,
		},
		&HistoryModelStruct{},
	)
}
