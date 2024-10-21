package controllers

import (
	"OAM/conf"
	"path"
)

type FileController struct {
	AuthController
}

// 上传图片，返回图片访问的接口url
func (c *FileController) UploadImg() {
	relpath, err := c.uploadFile("uploadfile", 0)
	if err != nil {
		c.JsonParamError(err.Error())
	} else {
		c.JsonOk(path.Join(conf.STATIC_EXT_BASE_URL, relpath))
	}
}

//上传普通文件
func (c *FileController) UploadFile() {
	relpath, err := c.uploadFile("uploadfile", 3)
	if err != nil {
		c.JsonParamError(err.Error())
	} else {
		c.JsonOk(path.Join(conf.STATIC_EXT_BASE_URL, relpath))
	}
}
