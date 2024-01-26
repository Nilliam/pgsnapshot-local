package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	for {

		file, err := os.Open("settings.json")
		if err != nil {
			log.Println("Error opening settings file:", err)
			panic(err)
		}
		defer file.Close()

		var settings Settings
		err = json.NewDecoder(file).Decode(&settings)
		if err != nil {
			log.Println("Error decoding settings file:", err)
			panic(err)
		}

		dump(settings)

		now := time.Now()

		var closestTime time.Time
		for _, backupTime := range settings.BackupTimes {
			parsedTime, err := time.Parse("15:04:05", backupTime)
			if err != nil {
				log.Println("Error parsing backup time:", err)
				return
			}
			if closestTime == (time.Time{}) || (parsedTime.After(now) && parsedTime.Before(closestTime)) {
				closestTime = parsedTime
			}
		}

		closestTime = time.Date(now.Year(), now.Month(), now.Day(), closestTime.Hour(), closestTime.Minute(), closestTime.Second(), 0, now.Location())

		// Calculate the duration until the next backup time
		now = time.Now()
		nextBackupTime := closestTime
		if nextBackupTime.Before(now) {
			nextBackupTime = nextBackupTime.Add(24 * time.Hour)
		}
		duration := nextBackupTime.Sub(now)

		// Sleep until the next backup time
		log.Printf("Next backup scheduled for %s", nextBackupTime.Format("2006-01-02 15:04:05"))
		time.Sleep(duration)

	}
}

func dump(settings Settings) {

	postgres := &pg.Postgres{
		Host:     settings.Connection.Host,
		Port:     settings.Connection.Port,
		Username: settings.Connection.User,
		Password: os.Getenv("PGPASSWORD"),
	}

	dump, err := pg.NewDump(postgres)
	if err != nil {
		panic(err)
	}

	dump.ResetOptions()
	plainText := "p"
	dump.SetupFormat(plainText)

	BackupFolders(settings)

	for _, database := range settings.Databases {
		dump.SetFileName(database + ".sql")

		postgres.DB = database

		dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: true, StreamDestination: os.Stdout})
		if dumpExec.Error != nil {
			fmt.Println("Dump failed")
			SendWebhook(settings.Webhook, "Dump failed")
			fmt.Println(dumpExec.Error.Err)
			fmt.Println(dumpExec.Output)
			panic(dumpExec.Error.Err)
		} else {
			fmt.Println("Dump success")
			fmt.Println(dumpExec.Output)
			zipDump := database + ".zip"
			AddFileToZip(zipDump, database+".sql")
			UploadToS3(settings, database+".zip")
		}
	}

	serverStatus(settings)
}

func serverStatus(settings Settings) {
	if settings.ServerUrl == "" {
		message := "The Backup Service is currently " +
			"running in local mode because no server configuration was found. " +
			"To enable synchronization, please add the server name and URL to the properties in settings.json."
		SendWebhook(settings.Webhook, message)
	}
}

func main() {

	serviceName := "PostgreSQL Backup Service"
	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceName,
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
