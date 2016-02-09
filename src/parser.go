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
	alias     string
	skip      bool
	fieldName string
	members   []interface{}
}

func squashSameStructs() {
	for i, s := range structList {
		for j := i + 1; j < len(structList); j++ {
			_struct := structList[j]
			checkAndModify(s, _struct)
		}
	}
}

func checkAndModify(struct1 *Struct, struct2 *Struct) {
	var fieldMap1 = make(map[string]string)
	var fieldMap2 = make(map[string]string)
	for _, f := range struct1.members {
		switch f.(type) {
		case *Struct:
			fieldMap1[f.(*Struct).fieldName] = f.(*Struct).alias
			continue
		}
		fieldMap1[f.(Field).Name] = f.(Field).Type
	}

	for _, f := range struct2.members {
		switch f.(type) {
		case *Struct:
			fieldMap2[f.(*Struct).fieldName] = f.(*Struct).alias
			continue
		}
		fieldMap2[f.(Field).Name] = f.(Field).Type
	}

	for k, v := range fieldMap1 {
		if v2, ok := fieldMap2[k]; !ok || v != v2 {
			return
		}
	}
	for k, v := range fieldMap2 {
		if v2, ok := fieldMap1[k]; !ok || v != v2 {
			return
		}
	}
	struct2.alias = struct1.alias
	struct2.skip = true
}

func exportStructs() {
	file, _ := os.OpenFile("./model.go", os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	for _, s := range structList {
		if s.skip {
			continue
		}
		file.Write([]byte("type " + s.fieldName + " struct {\n"))
		for _, f := range s.members {

			switch f.(type) {
			case *Struct:
				file.Write([]byte(fmt.Sprintf("\t%s\t%s\n", f.(*Struct).fieldName, f.(*Struct).alias)))
				continue
			}
			file.Write([]byte(fmt.Sprintf("\t%s\t%s\n", f.(Field).Name, f.(Field).Type)))
		}
		file.Write([]byte("}\n\n"))
	}
}

func handleMap(structIndex int, k string, v interface{}) {

	switch v.(type) {
	case map[string]interface{}:
		offset := 1
		if len(structList) >= structIndex+1 && len(structList) > 0 {
			offset = len(structList) - structIndex
		}
		structList = append(structList, &Struct{alias: k, fieldName: k})
		for _k, _v := range v.(map[string]interface{}) {
			handleMap(structIndex+offset, _k, _v)
		}
		if structIndex > -1 {
			structList[structIndex].members = append(structList[structIndex].members, structList[structIndex+offset])
		}
		break
	case nil:
		structList[structIndex].members = append(structList[structIndex].members, Field{k, "interface{}"})
		break
	default:
		structList[structIndex].members = append(structList[structIndex].members, Field{k, fmt.Sprintf("%T", v)})
	}
}

func main() {
	file, _ := os.Open("./test_files/real.json")
	var root interface{}
	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &root)
	handleMap(-1, "__root__", root)
	squashSameStructs()
	exportStructs()
}
