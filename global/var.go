package global

import "fmt"

var (
	fileDownloadApiPrefix = fmt.Sprintf("http://52.82.65.36:%d/%s", HttpPort, "api/file/download?name=%s")
)

func NewFileDownloadApiPath(filename string) string {
	return fmt.Sprintf(fileDownloadApiPrefix, filename)
}
