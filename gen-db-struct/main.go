package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

func toCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
func getSqlFieldNames(f *os.File) []string {
	var res []string

	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	for _, lineBytes := range bytes.Split(fBytes, []byte(`\n`)) {
		for _, wordBytes := range bytes.Split(lineBytes, []byte(`,`)) {
			res = append(res, string(bytes.TrimSpace(wordBytes)))
		}
	}

	return res
}
func namifyField(name string) string {
	return toCamelCase(name)
}

// last_name, first_name, middle_name,
//	suffix, entity_name, type,
//	gender, birth_date, death_date, tax_number,
//	created_on, created_by, updated_on, updated_by
func main() {
	structName := flag.String("struct", "Struct", "structure name")
	tagName := flag.String("tag", "gorm", "tag name to use")
	format := flag.String("format", "gorm", "Which format to use")
	flag.Parse()
	sqlFields := getSqlFieldNames(os.Stdin)
	fmt.Println("type", *structName, "struct {")
	for _, field := range sqlFields {
		if field == "" {
			continue
		}
		fieldGoName := namifyField(field)
		switch *format {
		case "gorm":
			fmt.Printf("\t %s string `%s:\"column:%s\"`\n", fieldGoName, *tagName, field)
		case "notag":
			fmt.Printf("\t %s\n", fieldGoName)
		default:
			panic("Unknown format " + *format)
		}
	}
	fmt.Println("}")
}
