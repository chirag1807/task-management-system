package repository

import (
	"reflect"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
)

func UpdateQuery(tableName string, model interface{}, id int64, flag int) (string, []interface{}, error) {
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

		if flag == 1 && field.Tag.Get("db") == "assignee_individual" {
			query += " " + "assignee_team" + " = $"
			query += strconv.Itoa(len(args)+1) + ","
			args = append(args, nil)
		} else if flag == 1 && field.Tag.Get("db") == "assignee_team" {
			query += " " + "assignee_individual" + " = $"
			query += strconv.Itoa(len(args)+1) + ","
			args = append(args, nil)
		}
	}
	query = query[:len(query)-1] + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)

	return query, args, nil
}

func UpdateTaskFields(dbTask response.Task, requestTask request.UpdateTask) response.Task {
	dbTaskType := reflect.TypeOf(dbTask)
	dbTaskValue := reflect.ValueOf(&dbTask).Elem()
	requestTaskValue := reflect.ValueOf(requestTask)

	for i := 0; i < dbTaskType.NumField(); i++ {
		field := dbTaskType.Field(i)
		fieldName := field.Name
		dbFieldValue := dbTaskValue.FieldByName(fieldName)
		requestFieldValue := requestTaskValue.FieldByName(fieldName)
		if requestFieldValue.IsValid() && !reflect.DeepEqual(requestFieldValue.Interface(), reflect.Zero(requestFieldValue.Type()).Interface()) {
			dbFieldValue.Set(requestFieldValue)
		}
	}

	return dbTask
}
