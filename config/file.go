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

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	if !exists(name) {
		return os.Create(name)
	}
	return os.OpenFile(name, flag, perm)
}
