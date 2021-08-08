package main

import (
	"srv/src/module/route"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Init database

	// Init routes
	r := route.Init()

	//r.Run("127.0.0.1:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
