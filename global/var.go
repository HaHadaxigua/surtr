package global

var (
	storagePath string
	apiAddr     string
)

func GetStoragePath() string {
	return storagePath
}

func setStoragePath(v string) {
	storagePath = v
}

func setApiAddr(v string) {
	apiAddr = v
}

func GetApiAddr() string {
	return apiAddr
}
