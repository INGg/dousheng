package util

import (
	"demo1/config"
	"fmt"
	"path/filepath"
)

func GetFileUrl(fileName string) string {
	base := fmt.Sprintf(`http://%s:%s/static/%s`, config.Info.IP, config.Info.Port, fileName)
	return base
}

func SaveImageFromVideo(fileName string, isDebug bool) error {
	v2i := NewVideo2Image()
	if isDebug {
		v2i.Debug()
	}
	v2i.InputPath = filepath.Join(config.Info.StaticSourcePath, fileName+defaultVideoSuffix)
	v2i.OutputPath = filepath.Join(config.Info.StaticSourcePath, fileName+defaultImageSuffix)
	v2i.FrameCount = 1
	queryString, err := v2i.GetQueryString()
	if err != nil {
		return err
	}
	return v2i.ExecCommand(queryString)
}
