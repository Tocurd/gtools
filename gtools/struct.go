package gtools

import (
	"errors"
	"reflect"
)

var Struct structInterface

type structInterface interface {
	Fields(value interface{}) ([]string, error)
	ToMap(data interface{}) (map[string]interface{}, error)
}

type selfStruct struct{}

func init() {
	Struct = &selfStruct{}
}

/**
 * @description: 获取结构体中的所有的字段名
 * @param {interface{}} value
 * @return {*}
 */
func (_selfStruct selfStruct) Fields(value interface{}) ([]string, error) {
	typeof := reflect.TypeOf(value)
	if typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
	}

	// 如果是map属性的话
	if typeof.Kind() == reflect.Map {
		valueof := reflect.ValueOf(value)
		if valueof.Kind() == reflect.Ptr {
			valueof = valueof.Elem()
		}
		result := make([]string, 0)
		maps := valueof.MapKeys()
		for index := 0; index < len(maps); index++ {
			result = append(result, maps[index].String())
		}
		return result, nil
	}

	// 强制要求必须输入结构体
	var fields = make([]string, 0)
	if typeof.Kind() != reflect.Struct {
		return fields, errors.New("check type error not struct")
	}

	for index := 0; index < typeof.NumField(); index++ {
		jsonKey := typeof.Field(index).Tag.Get("json")
		if jsonKey != "" {
			fields = append(fields, jsonKey)
		} else {
			fields = append(fields, typeof.Field(index).Name)
		}
	}

	return fields, nil
}

/**
 * @description: 结构体转map
 * @param {interface{}} data
 * @return {*}
 */
func (_selfStruct selfStruct) ToMap(data interface{}) (map[string]interface{}, error) {
	typeof := reflect.TypeOf(data)
	if typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
	}

	valueof := reflect.ValueOf(data)
	if valueof.Kind() == reflect.Ptr {
		valueof = valueof.Elem()
	}

	var result = make(map[string]interface{}, 1)
	// 强制要求必须输入结构体
	if typeof.Kind() != reflect.Struct {
		return result, errors.New("check type error not struct")
	}

	for index := 0; index < typeof.NumField(); index++ {
		jsonKey := typeof.Field(index).Tag.Get("json")

		fieldName := ""
		if jsonKey != "" {
			fieldName = jsonKey
		} else {
			fieldName = typeof.Field(index).Name
		}
		mapKeyInterface := valueof.FieldByName(typeof.Field(index).Name)
		if mapKeyInterface.Kind() == reflect.Invalid {
			return nil, errors.New("key not exist")
		}

		result[fieldName] = mapKeyInterface.Interface()

	}

	return result, nil
}
