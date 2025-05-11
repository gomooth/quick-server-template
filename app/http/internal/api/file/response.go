package file

type chunkResponse struct {
	Over bool   `json:"over"` // 分开上传文件是否全部上传完毕
	Url  string `json:"url"`  // Over = true 时，返回文件地址；否则为空
}
