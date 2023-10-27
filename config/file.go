package config

import "os"

// exists 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
	//获取文件信息
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
