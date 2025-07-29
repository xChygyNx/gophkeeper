package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func CreateStorageUser(dirPath string, id int64) error {
	userId := strconv.FormatInt(id, 10)
	path := filepath.Join(dirPath, userId)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error in mkdir %s: %w", path, err)
	}
	return nil
}

func CreateStorageNotExistsUser(dirPath string, id int64) error {
	userId := strconv.FormatInt(id, 10)
	path := filepath.Join(dirPath, userId)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error in mkdir %s: %w", path, err)
			}
		} else {
			return fmt.Errorf("error in check dir %s: %w", path, err)
		}
	}
	return nil
}

func UploadFile(dirPath string, id int64, name string, data []byte) error {
	userId := strconv.Itoa(int(id))
	path := filepath.Join(dirPath, userId, "/", name)
	// Write data to file
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("error in upload file %s: %w", path, err)
	}
	return nil
}

func DownloadFile(dirPath string, id int64, name string) ([]byte, error) {
	userId := strconv.FormatInt(id, 10)
	path := filepath.Join(dirPath, userId, "/", name)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error in open file %s: %w", path, err)
	}
	defer func() {
		err = file.Close()
	}()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error in read file %s: %w", path, err)
	}
	return data, nil
}

func RemoveFile(dirPath string, id int64, name string) error {
	userId := strconv.Itoa(int(id))
	path := filepath.Join(dirPath, userId, "/", name)
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("error in remove file %s: %w", path, err)
	}
	return nil
}
