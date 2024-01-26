package main

type Settings struct {
	ServerName string `json:"serverName"`
	ServerUrl  string `json:"serverUrl"`
	Webhook    string `json:"webhook"`
	S3Bucket   string `json:"s3Bucket"`
	AwsRegion  string `json:"awsRegion"`
	Connection struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
	} `json:"connection"`
	BackupTimes []string `json:"backupTimes"`
	Databases   []string `json:"databases"`
	Folders     []string `json:"folders"`
}
