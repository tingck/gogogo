package api

import (
	"github.com/chiti62/gogogo/ch2/internal/service"
	"github.com/chiti62/gogogo/ch2/pkg/app"
	"github.com/chiti62/gogogo/ch2/pkg/convert"
	"github.com/chiti62/gogogo/ch2/pkg/errcode"
	"github.com/chiti62/gogogo/ch2/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		logrus.WithContext(c).Errorf("svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
