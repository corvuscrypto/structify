package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var structList []*Struct

type Field struct {
	Name string
	Type string
}

type Struct struct {
	fieldName string
	members   []interface{}
}

// func calcNameDistance(name1 string, name2 string) float64 {
//   name1 = strings.
// }

func handleMap(structIndex int, k string, v interface{}) {
	fieldPrefix := ""

	//if not the first struct, then add the prefix
	if structIndex > -1 {
		fieldPrefix = structList[structIndex].fieldName + "."
	}

	switch v.(type) {
	case map[string]interface{}:
		structList = append(structList, &Struct{fieldName: fieldPrefix + k})
		for _k, _v := range v.(map[string]interface{}) {
			handleMap(structIndex+1, _k, _v)
		}
		if structIndex > -1 {
			structList[structIndex].members = append(structList[structIndex].members, structList[structIndex+1])
		}
		break
	case nil:
		structList[structIndex].members = append(structList[structIndex].members, Field{fieldPrefix + k, "interface{}"})
		break
	default:
		structList[structIndex].members = append(structList[structIndex].members, Field{fieldPrefix + k, fmt.Sprintf("%T", v)})
	}
}

func main() {
	file, _ := os.Open("./test_files/simple.json")
	var root interface{}
	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &root)
	handleMap(-1, "__root__", root)
	for _, v := range structList {
		fmt.Println(v)
	}
}
