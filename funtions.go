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

func ToFieldType(_type string) string {
	if _type == "field_type.DeletedTime" {
		return "int64"
	} else if _type == "time.Time" {
		return "int64"
	}
	return _type
}
