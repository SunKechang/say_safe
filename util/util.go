package util

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	SHORTID_DIGITS                = "abcdefghijkmnpqrstuvwxyz0123456789"
	DEFAULT_SHORTID_STRING        = "awsdaexy0yi1s02m"
	DEFAULT_PREFIX_SHORTID_STRING = "default-awsdaexy0yi1s02m"
)

func UUIDToShortID(UUID string) string {
	// 32uuid -> 32md5 hex
	data := []byte(UUID)
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash)

	var result []byte
	for i := 0; i < 16; i++ {
		// parse 2bit char from 16base to 10base
		index, _ := strconv.ParseUint(md5str[2*i:2*i+2], 16, 32)
		result = append(result, SHORTID_DIGITS[index%34])
	}
	return string(result)
}

func NewShortIDString(prefix string) string {
	needPrefix := prefix != ""
	shortidStr := DEFAULT_SHORTID_STRING
	if needPrefix {
		shortidStr = DEFAULT_PREFIX_SHORTID_STRING
	}

	newID := uuid.NewV4()
	shortidStr = UUIDToShortID(newID.String())
	if needPrefix {
		shortidStr = prefix + "-" + shortidStr
	}
	return shortidStr
}

func ZipFile(sourceFilePath, targetPath, password string) error {
	strs := strings.Split(sourceFilePath, "/")
	fileName := strs[len(strs)-1]
	fileDir := sourceFilePath[0 : len(sourceFilePath)-len(fileName)]

	var command string
	if password != "" {
		command = fmt.Sprintf("cd %s && zip -r -P %s %s %s", fileDir, password, targetPath, fileName)
	} else {
		command = fmt.Sprintf("cd %s && zip -r %s %s", fileDir, targetPath, fileName)
	}

	cmd := exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		fmt.Printf("Zip file[%s] command[%s] error: %v", sourceFilePath, command, err)
		return err
	}
	return nil
}

func SplitFile(filePath, singleSize string) ([]string, error) {
	strs := strings.Split(filePath, "/")
	fileName := strs[len(strs)-1]
	fileDir := filePath[0 : len(filePath)-len(fileName)]
	command := fmt.Sprintf("cd %s && split -b %s -d %s %s", fileDir, singleSize, fileName, fileName)
	cmd := exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		fmt.Printf("Split file[%s] command[%s] error: %v", filePath, command, err)
		return nil, err
	}

	command2 := fmt.Sprintf("cd %s && ls -r | grep %s", fileDir, fileName)
	cmd2 := exec.Command("/bin/bash", "-c", command2)
	out, err := cmd2.Output()
	if err != nil {
		fmt.Printf("Get splited files[%s] command[%s] error: %v", filePath, command2, err)
		return nil, err
	}
	files := strings.Split(string(out), "\n")
	length := len(files)
	for i := 0; i < length; i++ {
		if files[i] == fileName || files[i] == "" {
			files = append(files[0:i], files[i+1:]...)
			i--
			length--
		} else {
			files[i] = path.Join(fileDir, files[i])
		}

	}
	return files, nil
}
