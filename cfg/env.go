package cfg

type Config struct {
	DBUsername string `envconfig:"DB_USER"`
	DBPassword string `envconfig:"DB_PASS"`
	DBHost     string `envconfig:"DB_HOST"`
	DBPort     int    `envconfig:"DB_PORT"`
	DBName     string `envconfig:"DB_NAME"`

	//StorageDriver                 string `envconfig:"STORAGE_DRIVER" default:"local"`
	//StoragePath                   string `envconfig:"STORAGE_PATH"`
	//LocalStorageDownloadPrefixUrl string `envconfig:"LOCAL_STORAGE_DOWNLOAD_PREFIX_URL"`
	//
	//AWSProfile   string `envconfig:"AWS_PROFILE"`
	S3BucketName string `envconfig:"S3_BUCKET_NAME"`
}
