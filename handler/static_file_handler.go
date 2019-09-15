package handler

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"github.com/gingerxman/eel/utils"
	"github.com/gingerxman/eel/log"
)

func HandleStaticFile(path string, response *Response) bool {
	curDir, _ := filepath.Abs(filepath.Dir("."))
	
	absPath := filepath.Join(curDir, path)
	isExists := utils.FileExists(absPath)
	log.Logger.Debugw("check static file", "path", absPath, "exists", isExists)
	if isExists {
		file, err := os.Open(absPath)
		if err != nil {
			return false
		}
		defer file.Close()
		
		var bufferWriter bytes.Buffer
		io.Copy(&bufferWriter, file)
		
		contentType := "text/html; charset=utf-8"
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "text/javascript"
		}
		response.Header("Content-Type", contentType)
		response.Body(bufferWriter.Bytes())
		
		return true
	}
	
	return false
}
