package services

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	cacheExpireTime = 15
)

func getCacheDirPath() (*string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	dir := homedir + "/.config/eccu"
	return &dir, nil
}

func getCache() (*string, error) {
	cacheDirPath, err := getCacheDirPath()
	if err != nil {
		return nil, err
	}

	cacheFilePath := *cacheDirPath + "/.cache"
	result, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	var s syscall.Stat_t
	syscall.Stat(cacheFilePath, &s)
	now := time.Now()

	sec, _ := s.Mtim.Unix()
	lastFileModifiedTime := time.Unix(sec, 0)

	duration := now.Sub(lastFileModifiedTime)
	log.Printf("[DEBUG] %s", duration)
	if duration.Minutes() > cacheExpireTime {
		return nil, errors.New("cache expired")
	}
	data := string(result)
	return &data, nil
}

func putCache(data string) error {
	cacheDirPath, err := getCacheDirPath()
	if err != nil {
		return err
	}

	if f, err := os.Stat(*cacheDirPath); os.IsNotExist(err) || f.IsDir() {
		err := os.MkdirAll(*cacheDirPath, 0777)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(*cacheDirPath + "/.cache")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(([]byte)(data))
	return nil
}
