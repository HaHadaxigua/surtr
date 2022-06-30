package global

import "fmt"

var (
	fileDownloadApiPrefix = "http://52.82.65.36/api/file/download?name=%s"
)

func NewFileDownloadApiPath(filename string) string {
	return fmt.Sprintf(fileDownloadApiPrefix, filename)
}
