package controller

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	im_toolkit "github.com/zengyu2020/im-toolkit"
	"go.uber.org/zap"
)

const UploadFilePath = "./upload/static"

var (
	ErrInternalError      = im_toolkit.NewCustomError(100, "服务器内部发生错误")
	ErrPermissionDenied   = im_toolkit.NewCustomError(101, "没有上传的权限")
	ErrFileSizeOutOfLimit = im_toolkit.NewCustomError(102, "文件大小超出限制，限10M以内的文件")
	ErrFileTypeInvalid    = im_toolkit.NewCustomError(103, "不合法的文件格式")
	ErrInvalidArguments   = im_toolkit.NewCustomError(104, "错误的请求参数")
	ErrFileNotExist       = im_toolkit.NewCustomError(105, "文件不存在")
)

type StaticFileController struct {
}

func NewStaticFileController() *StaticFileController {
	return &StaticFileController{}
}

func (ctrl *StaticFileController) Get(redisC redis.UniversalClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := new(GetReq)
		if err := ctx.Bind(req); err != nil {
			zap.L().Error("绑定query发生错误", zap.Error(err))
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// if len(req.Token) == 0 {
		// 	zap.L().Error("请求中没有token")
		// 	ctx.AbortWithStatus(http.StatusBadRequest)
		// 	return
		// }

		// if req.Token != "JKKDLSWNSMSAIIEHHHDNX" {

		// 	// 验证是否有权限上传
		// 	authorization := req.Token
		// 	count, err := redisC.Exists(fmt.Sprintf("upload_token:%s", authorization)).Result()
		// 	if err != nil {
		// 		zap.L().Error("向redis查询upload_token发生错误", zap.Error(err))
		// 		ctx.AbortWithStatus(http.StatusBadRequest)
		// 		return
		// 	}

		// 	if count != 1 {
		// 		zap.L().Error("redis不存在token", zap.String("token", authorization))
		// 		ctx.AbortWithStatus(http.StatusBadRequest)
		// 		return
		// 	}
		// }

		_, err := uuid.Parse(strings.Split(req.Id, ".")[0])
		if err != nil {
			zap.L().Error("id不是uuid", zap.Error(err))
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, err = os.Stat(filepath.Join(UploadFilePath, req.Id))
		if errors.Is(err, os.ErrNotExist) {
			zap.L().Error("请求的文件不存在", zap.Error(err))
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		} else if err != nil {
			zap.L().Error("查找文件时发生错误", zap.Error(err))
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.File(filepath.Join(UploadFilePath, req.Id))
	}
}

func (ctrl *StaticFileController) Upload(redisC redis.UniversalClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 验证是否有权限上传
		// mf, err := ctx.MultipartForm()
		// if err != nil {
		// 	zap.L().Error("转换提交表单时发生错误", zap.Error(err))
		// 	im_toolkit.Response(ctx, ErrInternalError)
		// 	return
		// }

		// appSecret, ok := mf.Value["appSecret"]
		// if !ok {
		// 	im_toolkit.Response(ctx, ErrInvalidArguments)
		// 	return
		// }

		// if len(appSecret) != 1 {
		// 	im_toolkit.Response(ctx, ErrInvalidArguments)
		// 	return
		// }

		// authorization := appSecret[0]

		// if len(authorization) == 0 {
		// 	zap.L().Error("请求中没有token")
		// 	im_toolkit.Response(ctx, ErrInvalidArguments)
		// 	return
		// }

		// if authorization != "JKKDLSWNSMSAIIEHHHDNX" {
		// 	count, err := redisC.Exists(fmt.Sprintf("upload_token:%s", authorization)).Result()
		// 	if err != nil {
		// 		zap.L().Error("向redis查询upload_token发生错误", zap.Error(err))
		// 		im_toolkit.Response(ctx, ErrInvalidArguments)
		// 		return
		// 	}

		// 	if count != 1 {
		// 		zap.L().Error("redis不存在token", zap.String("token", authorization))
		// 		im_toolkit.Response(ctx, ErrInvalidArguments)
		// 		return
		// 	}
		// }

		header, err := ctx.FormFile("file")
		if err != nil {
			zap.L().Error("获取上传文件时发生错误", zap.Error(err))
			im_toolkit.Response(ctx, ErrInternalError)
			return
		}
		if header.Size > 1024*1024*10 {
			zap.L().Error("错误,图片太大,上传失败")
			im_toolkit.Response(ctx, ErrFileSizeOutOfLimit)
			return
		}

		f, err := header.Open()
		if err != nil {
			zap.L().Error("打开上传文件时发生错误", zap.Error(err))
			im_toolkit.Response(ctx, ErrInternalError)
			return
		}

		defer func() {
			if err := f.Close(); err != nil {
				zap.L().Error("关闭上传文件时发生错误", zap.Error(err))
			}
		}()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			zap.L().Error("读取上传文件全部字节时发生错误", zap.Error(err))
			im_toolkit.Response(ctx, ErrInternalError)
			return
		}

		kind, err := filetype.Image(buf)
		if err != nil {
			zap.L().Error("确认上传文件类型时发生错误", zap.Error(err))
			im_toolkit.Response(ctx, ErrInternalError)
			return
		}

		switch kind {
		case matchers.TypeJpeg, matchers.TypePng:
			filename := uuid.New().String() + "." + kind.Extension

			if err := ioutil.WriteFile(filepath.Join(UploadFilePath, filename), buf, os.ModePerm); err != nil {
				zap.L().Error("保存上传文件时发生错误", zap.Error(err))
				im_toolkit.Response(ctx, ErrInternalError)
				return
			}

			im_toolkit.Response(ctx, filename)
		default:
			im_toolkit.Response(ctx, ErrFileTypeInvalid)
			return
		}
	}
}

//
//func (ctrl *StaticFileController) Delete(ctx *gin.Context) {
//	fileId := ctx.Param("id")
//	if len(fileId) == 0 {
//		ctx.AbortWithStatus(http.StatusNotFound)
//		return
//	}
//
//	err := os.Remove(filepath.Join(UploadFilePath, fileId))
//	if Is(err, os.ErrPermission) {
//		zap.L().Error("删除文件时权限不够", zap.Error(err))
//		ctx.AbortWithStatus(http.StatusInternalServerError)
//		return
//	} else if !Is(err, os.ErrNotExist) {
//		zap.L().Error("删除文件时发生错误", zap.Error(err))
//		ctx.AbortWithStatus(http.StatusInternalServerError)
//		return
//	}
//}
