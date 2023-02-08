package main

import (
	"smart_home/ctrl"
	"smart_home/db"
)

func main() {
	database := db.NewDatabase()
	control := ctrl.MustNewService(database)
	r := control.SetRouter()
	err := r.Run(":8081")
	if err != nil {
		panic("Run Engine Failure!")
	}
}
