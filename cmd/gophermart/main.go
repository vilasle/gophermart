package main

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
)

func main() {
	sb := sqlbuilder.Update("order")
	sb.Set(fmt.Sprintf("status = %d", 3))
	sb.Where(sb.Equal("number", "213123"))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_ = args

	fmt.Println(txt)
}
