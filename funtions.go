package gen

import "gorm.io/gen/internal/model"

func Add(a, b int) int {
	return a + b
}

func ExistsField(fieldName string, fields []*model.Field) bool {
	for _, field := range fields {
		if field.Name == fieldName {
			return true
		}
	}
	return false
}
