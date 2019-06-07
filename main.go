package main

import (
	"github.com/jinzhu/gorm"
	"github.com/novacloudcz/graphql-orm/cmd"
)

var db *gorm.DB

func main() {
	cmd.Execute()
}
