package utils

import (
	"encoding/json"
)

func ParseIDs2String(ids []int64) (string, error) {
	bytes, err := json.Marshal(ids)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseIDString2IDs(str string) ([]int64, error) {
	ids := make([]int64, 0)
	err := json.Unmarshal([]byte(str), &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func RemoveRepeatedIDs(ids []int64) []int64 {
	uintMap := make(map[int64]int64)
	for _, v := range ids {
		uintMap[v] = v
	}
	var res []int64
	for _, value := range uintMap {
		res = append(res, value)
	}
	return res
}
