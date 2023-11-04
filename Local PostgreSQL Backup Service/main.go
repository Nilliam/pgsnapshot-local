package main

import (
	"fmt"
	"log"

	"os"

	pg "github.com/habx/pg-commands"
	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func (p *program) run() {
	dump, err := pg.NewDump(&pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "database",
		Username: "postgres",
		Password: os.Getenv("PGPASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	dump.ResetOptions()
	plainText := "p"
	dump.SetupFormat(plainText)
	dump.SetFileName("/Preferred_folder/dump.sql")

	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: true, StreamDestination: os.Stdout})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)

	} else {
		fmt.Println("Dump success")
		fmt.Println(dumpExec.Output)
	}
}

func main() {
	svcConfig := &service.Config{
		Name:        "PostgreSQL Backup Service",
		DisplayName: "PostgreSQL Backup Service",
		Description: "PostgreSQL Backup Service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
