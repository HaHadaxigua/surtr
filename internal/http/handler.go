package http

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/HaHadaxigua/surtr/internal/http/file"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func routers(r gin.IRouter) {
	fileGroup := r.Group("/file")
	fileGroup.GET("/list", list)
	fileGroup.GET("/download", download)
	fileGroup.POST("/upload", upload)
	fileGroup.GET("/sleep", sleep)
}

// List
// @Tag File API
// @Title List File
// @Description list files
// @Success 200 {object} file.ListResp The Response object
// @router /api/file/list [get]
// list will return files on current machine
func list(c *gin.Context) {
	resp, err := file.New().List()
	if err != nil {
		logrus.Errorf("failed to list files: %v", err)
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	c.JSON(http.StatusOK, Ok(resp))
}

// download will download file with expected file name
func download(c *gin.Context) {
	var req file.DownloadReq
	if err := c.BindQuery(&req); err != nil {
		logrus.Errorf("failed to parse query for %v: %v", reflect.TypeOf(req), err)
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	resp, err := file.New().Download(&req)
	if err != nil {
		logrus.Errorf("failed to download: %s : %v", req.Filename, err)
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.Filename))
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", resp.Size)
	c.Writer.Write(resp.Data)
}

// upload will upload file to the fixed path
func upload(c *gin.Context) {
	fh, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	var req file.UploadReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}
	req.FileHeader = header
	req.File = fh

	if err := file.New().Upload(&req); err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}
	c.JSON(http.StatusOK, Ok(nil))
}

func sleep(c *gin.Context) {
	type req struct {
		timeS int `json:"time" form:"time"`
	}
	var r req
	if err := c.BindQuery(&r); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	time.Sleep(time.Duration(r.timeS) * time.Second)
	c.JSON(http.StatusOK, Ok("hello"))
}
