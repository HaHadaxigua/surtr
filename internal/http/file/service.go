package file

import (
	"errors"
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
	"github.com/sirupsen/logrus"
)

type service struct {
	storage string
}

func New() *service {
	return &service{
		storage: global.Conf.Storage,
	}
}

type DownloadReq struct {
	Filename string `form:"name" json:"name" binding:"required"`
}

type DownloadResp struct {
	Info

	Data []byte
}

func (s *service) Download(req *DownloadReq) (resp *DownloadResp, err error) {
	lookupPath := filepath.Join(s.storage, req.Filename)
	if util.IsNotExist(lookupPath) {
		return nil, fmt.Errorf("cannot find expected file: %s", lookupPath)
	}

	info, err := newInfo(lookupPath)
	if err != nil {
		return nil, err
	}

	if info.IsDir {
		return nil, errors.New("not support download folder at present")
	}

	buf, err := os.ReadFile(lookupPath)
	if err != nil {
		return nil, err
	}

	return &DownloadResp{
		Info: *info,
		Data: buf,
	}, nil
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

type ListResp struct {
	List []*ListFileItem
}

type ListFileItem struct {
	Info
	Link string
}

func (s *service) List() (*ListResp, error) {
	dirEntries, err := os.ReadDir(s.storage)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder: %s, %v", s.storage, err)
	}

	var resp ListResp
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			info, err := newInfo(entry.Name())
			if err != nil {
				logrus.Infof("failed to open file: %s, %v", entry.Name(), err)
				return nil, err
			}
			resp.List = append(resp.List, &ListFileItem{
				Info: *info,
				Link: global.NewFileDownloadApiPath(entry.Name()),
			})
		}
	}

	return &resp, nil
}

type Info struct {
	Filename   string `json:"filename"`
	IsDir      bool   `json:"isDir,omitempty"`
	Size       string `json:"size,omitempty"`
	Permission string `json:"permission"`
	ModifyTime string `json:"modifyTime,omitempty"`
}

func newInfo(filename string) (*Info, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load info of %s: %v", filename, err)
	}
	t, err := times.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load time info of %s: %v", filename, err)
	}

	info := &Info{
		Filename:   fi.Name(),
		Permission: fi.Mode().String(),
		ModifyTime: carbon.Time2Carbon(t.ModTime()).String(),
	}
	if util.IsFileExist(filename) {
		info.Size = humanize.Bytes(uint64(fi.Size()))
	} else {
		info.IsDir = true
	}

	return info, nil
}
