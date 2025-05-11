package file

import (
	"mime/multipart"
	"path"
	"strings"

	"github.com/save95/xerror"
)

type uploadRequest struct {
	Genre    string // 文件类型
	Business string // 业务类型

	Params map[string]interface{} // 业务参数

	File *multipart.FileHeader
}

func (in *uploadRequest) Validate() error {
	if in.File == nil {
		return xerror.New("请选择上传文件")
	}

	ext := strings.ToLower(path.Ext(in.File.Filename))

	switch in.Genre {
	case "pictures":
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" {
			return xerror.New("图片格式不对")
		}
	case "videos":
		if ext != ".mp4" {
			return xerror.New("视频格式不对")
		}
	case "attachments":
		if ext != ".docx" && ext != ".doc" &&
			ext != ".xlsx" && ext != ".xls" &&
			ext != ".pptx" && ext != ".ppt" &&
			ext != ".pdf" {
			return xerror.New("附件仅支持：word, excel, ppt, pdf")
		}
	default:
		return xerror.New("不支持的上传类型")
	}

	return nil
}

type base64Request struct {
	Genre    string // 文件类型
	Business string // 业务类型

	Filename string                 `json:"filename"` // 文件名
	Format   string                 `json:"format"`   // 文件格式
	Params   map[string]interface{} `json:"params"`   // 业务参数
	Data     string                 `json:"data"`     // Base64 内容
}

func (in *base64Request) Validate() error {
	if len(in.Filename) == 0 {
		return xerror.New("文件名 不能为空")
	}
	if len(in.Format) == 0 {
		return xerror.New("文件格式 不能为空")
	}
	if len(in.Data) == 0 {
		return xerror.New("base64 内容不能为空")
	}

	switch in.Genre {
	case "pictures":
		if !strings.HasPrefix(in.Format, "image/") {
			return xerror.New("图片格式不对")
		}
	default:
		return xerror.New("不支持的上传类型")
	}

	return nil
}

type chunkRequest struct {
	Genre    string // 文件类型
	Business string // 业务类型

	Params map[string]interface{} // 业务参数

	File *multipart.FileHeader

	FileId      int
	ChunksTotal int
	ChunkNumber int
}

func (in *chunkRequest) Validate() error {
	if in.File == nil {
		return xerror.New("请选择上传文件")
	}

	if in.FileId == 0 || in.ChunksTotal == 0 || in.ChunkNumber == 0 {
		return xerror.New("分块参数错误")
	}

	return nil
}
