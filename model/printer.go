package model

import (
	"github.com/graphql-go/graphql/language/printer"
)

// PrintSchema
func PrintSchema(model Model) (string, error) {

	printed := printer.Print(model.Doc)
	printedString, _ := printed.(string)
	// if ok {
	// 	schema += printedString + "\n"
	// 	fmt.Println("??", printed)
	// }
	return printedString, nil
}
