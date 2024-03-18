package repository

import (
	"fmt"
	"reflect"
	"strconv"
)

func UpdateQuery(tableName string, model interface{}, id int64) (string, []interface{}, error) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	query := "UPDATE " + tableName + " SET"
	var args []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		if field.Tag.Get("json") == "id" {
			continue
		}

		if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
			continue
		}

		query += " " + field.Tag.Get("db") + " = $"
		query += strconv.Itoa(len(args)+1) + ","
		args = append(args, fieldValue.Interface())
	}

	query = query[:len(query)-1] + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)

	fmt.Println(query, args)

	return query, args, nil
}
