package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "gopkg.in/clog.v1"
)

// Unzip file `srcZip` to given folder `destFolder`
//
func Unzip(srcZip string, destFolder string) error {
	if _, err := os.Stat(srcZip); os.IsNotExist(err) {
		return fmt.Errorf("input is not an zip file")
	}
	if st, err := os.Stat(destFolder); os.IsNotExist(err) {
		err = os.MkdirAll(destFolder, 666)
		if err != nil {
			return fmt.Errorf("failed to create destination folder")
		}
	} else if !st.IsDir() {
		return fmt.Errorf("destination is not an valid folder")
	}

	log.Info("[Unzip] Unzip file %s.", srcZip)
	start := time.Now()

	rzip, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer rzip.Close()

	for _, file := range rzip.File {
		var (
			err error
			fd  *os.File
			fc  io.ReadCloser
		)

		fpath := filepath.Join(destFolder, file.Name)
		//log.Trace("[Unzip] file : %s.", file.Name)
		//fmt.Printf("[Unzip] file : %s\n", file.Name)

		if file.FileInfo().IsDir() {
			os.Mkdir(fpath, file.Mode())
			if err != nil {
				log.Error(2, "[Unzip] Create dir %s failed with %v.", fpath, err)
			}
			continue
		} else {
			idx := strings.LastIndex(fpath, string(os.PathSeparator))
			if idx == -1 {
				log.Error(2, "[Unzip] Invalid file name %s.", fpath)
				continue
			}
			ppath := fpath[:idx]
			if err = os.MkdirAll(ppath, file.Mode()); err != nil {
				log.Error(2, "[Unzip] Create folder %s failed with %v.", err)
				continue
			}
		}

		for {
			if fc, err = file.Open(); err != nil {
				log.Error(2, "[Unzip] Open zip file %s failed with %v.", file.Name, err)
				break
			}
			if fd, err = os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode()); err != nil {
				log.Error(2, "[Unzip] Create file %s failed with %v.", fpath, err)
				break
			}
			if _, err = io.Copy(fd, fc); err != nil {
				log.Error(2, "[Unzip] Copy file failed with %v.", err)
				break
			}
			break
		}
		if fc != nil {
			fc.Close()
		}
		if fd != nil {
			fd.Close()
		}
		if err != nil {
			return err
		}
	}

	log.Info("[Unzip] Cost %s.", time.Since(start))
	return nil
}
