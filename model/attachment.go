package model

import (
	"Etpmls-Admin-Server/database"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type Attachment struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Path string	`json:"path"`
	OwnerID uint	`json:"owner_id"`
	OwnerType string	`json:"owner_type"`
}

// Validate if file is a image
// 验证文件是否为图片
func AttachmentValidateImageV2(h *multipart.FileHeader) (s string, err error) {
	f, err := h.Open()
	if err != nil {
		return s, err
	}

	// 识别图片类型
	_, image_type, _ := image.Decode(f)

	// 获取图片的类型
	switch image_type {
	case `jpeg`:
		return "jpeg", nil
	case `png`:
		return "png", nil
	case `gif`:
		return "git", nil
	case `bmp`:
		return "bmp", nil
	default:
		err := errors.New("This is not an image file, or the image file format is not supported!")
		return "", err
	}
}

// Upload Image
// 上传图片
func AttachmentUploadImageV2(c *gin.Context, file *multipart.FileHeader, extension string) (p string, err error) {
	// Make Dir
	t := time.Now().Format("20060102")
	path := "storage/upload/" + t + "/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return p, err
	}
	// UUID File name
	u := strings.ReplaceAll(uuid.New().String(), "-", "")

	file_path := path + u + "." + extension
	err = c.SaveUploadedFile(file, file_path)
	if err != nil {
		return p, errors.New("保存文件失败！")
	}

	if err = database.DB.Create(&Attachment{Path: file_path}).Error; err != nil {
		return p, err
	}

	return file_path, err
}

// Validate Path is a file in storage/upload
// 严重路径是否在storage/upload中
func AttachmentValidatePathV2(path string) error {
	const upload_path = "storage/upload/"
	// 截取前十五个字符判断和Path是否相同
	if len(path) <= len(upload_path) || !strings.Contains(path[:len(upload_path)], upload_path) {
		return  errors.New("请求路径非法！")
	}
	// 删除前缀
	// p := strings.TrimPrefix(path, upload_path)
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// 如果文件是目录
	if f.IsDir() {
		return errors.New("无法删除目录！")
	}
	return nil
}

type ApiAttachmentDeleteImageV2 struct {
	ID uint `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Path string	`json:"path" binding:"required"`
}
// Delete Image
// 删除图片
func AttachmentDeleteImageV2(j ApiAttachmentDeleteImageV2) (err error) {
	err = os.Remove(j.Path)
	if err != nil {
		return err
	}

	// Delete Database
	if err = database.DB.Unscoped().Where("path = ?", j.Path).Delete(Attachment{}).Error; err != nil {
		return err
	}

	return err
}

// Batch delete any type of files in storage/upload/
// 批量删除storage/upload/中的任何类型文件
func AttachmentBatchDeleteV2(s []string) (err error) {
	for _, v := range s {
		// Validate If a File
		err = AttachmentValidatePathV2(v)
		if err != nil {
			return err
		}
		// Delete Image
		_ = os.Remove(v)
	}

	// Delete Database
	if err = database.DB.Unscoped().Where("path IN (?)", s).Delete(Attachment{}).Error; err != nil {
		return err
	}

	return err
}


