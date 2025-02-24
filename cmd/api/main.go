package main

import "ambic/internal/bootstrap"

func main() {
	if err := bootstrap.Start(); err != nil {
		panic(err)
	}
}
