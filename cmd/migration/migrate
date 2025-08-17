package main

import "internal/infra/mysql"

func main() {
	if err := mysql.Migrate(); err != nil {
		panic(err)
	}
}
