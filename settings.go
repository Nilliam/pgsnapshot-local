package main

type Settings struct {
	ServerName                string `json:"serverName"`
	ServerUrl                 string `json:"serverUrl"`
	Webhook                   string `json:"webhook"`
	S3Bucket                  string `json:"s3Bucket"`
	S3Directory               string `json:"s3Directory"`
	AwsAccessKeyIdEncoded     string `json:"awsAccessKeyIdEncoded"`
	AwsSecretAccessKeyEncoded string `json:"awsSecretAccessKeyEncoded"`
	AwsRegion                 string `json:"awsRegion"`
	Connection                struct {
		Host            string `json:"host"`
		Port            int    `json:"port"`
		User            string `json:"user"`
		PasswordEncoded string `json:"passwordEncoded"`
	} `json:"connection"`
	BackupTimes []string `json:"backupTimes"`
	Databases   []string `json:"databases"`
	Folders     []string `json:"folders"`
}
