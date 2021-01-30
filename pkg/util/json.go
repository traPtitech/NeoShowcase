package util

import "encoding/json"

// ToJSON vをJSONに変換
// vは必ずJSONに変換できる型
// 失敗した場合はpanicします
func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// FromJSON strをmap[string]interface{}に変換
// strは必ず正しいJSON文字列
// 失敗した場合はpanicします
func FromJSON(str string) (v map[string]interface{}) {
	if err := json.Unmarshal([]byte(str), &v); err != nil {
		panic(err)
	}
	return v
}
