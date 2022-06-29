package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/HaHadaxigua/surtr/global"
	"github.com/HaHadaxigua/surtr/util"

	"github.com/djherbis/times"
	"github.com/dustin/go-humanize"
	"github.com/golang-module/carbon/v2"
)

type service struct {
	storage string
}

func New(v ...string) *service {
	storage := global.StoragePath
	if v != nil {
		storage = v[0]
	}
	return &service{
		storage: storage,
	}
}

type DownloadReq struct {
	Filename string `form:"name" json:"name" binding:"required"`
}

type DownloadResp struct {
	Filename   string `json:"filename"`
	IsDir      bool   `json:"isDir,omitempty"`
	Size       string `json:"size,omitempty"`
	Permission string `json:"permission"`
	ModifyTime string `json:"modifyTime,omitempty"`
	Data       []byte `json:"data,omitempty"`
}

func (s *service) Download(req *DownloadReq) (resp *DownloadResp, err error) {
	lookupPath := filepath.Join(s.storage, req.Filename)
	if util.IsNotExist(lookupPath) {
		return nil, fmt.Errorf("cannot find expected file: %s", lookupPath)
	}

	fi, err := os.Stat(lookupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load info of %s: %v", lookupPath, err)
	}
	t, err := times.Stat(lookupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load time info of %s: %v", lookupPath, err)
	}

	resp = &DownloadResp{
		Filename:   fi.Name(),
		Permission: fi.Mode().String(),
		ModifyTime: carbon.Time2Carbon(t.ModTime()).String(),
	}

	if util.IsFileExist(lookupPath) {
		resp.Size = humanize.Bytes(uint64(fi.Size()))
	} else {
		resp.IsDir = true
		return resp, fmt.Errorf("not support download directory: %s", lookupPath)
	}

	return
}

type UploadReq struct {
	File       io.Reader             `json:"file"`
	FileHeader *multipart.FileHeader `json:"fileHeader"`
}

func (s *service) Upload(req *UploadReq) error {
	buf, err := io.ReadAll(req.File)
	if err != nil {
		return err
	}
	return util.CreateFile(filepath.Join(s.storage, req.FileHeader.Filename), buf)
}
