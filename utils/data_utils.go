package utils

import (
	"encoding/json"
	"errors"
	"sort"

	"casorder/db/models"
	"casorder/utils/types"
)

func ValidateArgs(obj types.JSON, requiredParams []string) error {
	for _, val := range requiredParams {
		if obj[val] == nil {
			msg := "Parameter " + val + " is required"
			return errors.New(msg)
		}
	}
	return nil
}

func ParseObjectFromJson(jsonObj interface{}, model interface{}) error {
	jsonByte, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonByte, &model); err != nil {
		return err
	}

	return nil
}


func Reduce(f func(x1 models.OrderProduct, x2 models.OrderProduct) models.OrderProduct,
	iter []models.OrderProduct, initial *models.OrderProduct)  models.OrderProduct{
	var value models.OrderProduct
	if initial == nil {
		value = iter[0]
		iter = iter[1:]
	} else {
		value = *initial
	}
	if f == nil {
		f = SumFunc
	}
	for _, obj := range iter {
		value = f(value, obj)
	}
	return value
}

func SumFunc(x1 models.OrderProduct, x2 models.OrderProduct) models.OrderProduct{
	quantity := x1.Quantity + x2.Quantity
	x1.Quantity = quantity

	var data1 types.JSON
	var data2 types.JSON
	_ = json.Unmarshal(x1.Data, &data1)
	_ = json.Unmarshal(x2.Data, &data2)

	if data1["volumes"] != nil || data2["volumes"] != nil {
		data1["volumes"] = append(data1["volumes"].([]interface{}), data2["volumes"].([]interface{})...)
	}

	jsonByte, _ := json.Marshal(data1)
	_ = json.Unmarshal(jsonByte, &x1.Data)

	return x1
}


func Group(iter []models.OrderProduct) [][]models.OrderProduct{
	sort.SliceStable(iter, func(i, j int) bool {
		return iter[i].ProductID < iter[j].ProductID
	})

	var result [][]models.OrderProduct
	tmp := iter[0]
	group := make([]models.OrderProduct, 0)
	for _, v := range iter {
		if tmp.ProductID != v.ProductID {
			result = append(result, group)
			group = make([]models.OrderProduct, 0)
			group = append(group, v)
			tmp = v
		}else{
			group = append(group, v)
		}
	}
	result = append(result, group)
	return result
}

func Contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
