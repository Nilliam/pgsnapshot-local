package main

type Settings struct {
	Connection struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
	} `json:"connection"`
	BackupTimes []string `json:"backupTimes"`
	Databases   []string `json:"databases"`
	Folders     []string `json:"folders"`
}
