package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Response1 struct {
	Page   int
	Fruits []string
}

type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {
	// Encoding: go object -> json string
	res1_data := Response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "banana"}}
	res1_byte, _ := json.Marshal(res1_data)
	fmt.Println(string(res1_byte))

	res2_data := Response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "banana"}}
	res2_byte, _ := json.Marshal(res2_data)
	fmt.Println(string(res2_byte))

	// Decoding: json string -> go object
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat map[string]interface{} // necessary to map string to any data structure

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	fmt.Printf("type of num: %T\n", dat["num"]) // float64
	num := dat["num"].(float64)
	fmt.Println(num)

	fmt.Printf("type of strs: %T\n", dat["strs"]) // []interface{}
	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	// Decoding: with predifined struct
	res_byte := []byte(`{"page": 7,"fruits":["apple","banana"]}`)
	res_data := Response2{}
	if err := json.Unmarshal(res_byte, &res_data); err != nil {
		panic(err)
	}
	fmt.Println(res_data)

	// Encoder and Decoder io stream
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(res_data) // Direct encode an object and write its json to stdout
}
