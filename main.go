package main

import (
	"fmt"

	"os"

	pg "github.com/habx/pg-commands"
)

func main() {
	dump, err := pg.NewDump(&pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "",
		Username: "postgres",
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	dump.ResetOptions()
	dump.SetupFormat("p")
	dump.SetFileName("/local_dir/dump.sql")

	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: true, StreamDestination: os.Stdout})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)

	} else {
		fmt.Println("Dump success")
		fmt.Println(dumpExec.Output)
	}
}
