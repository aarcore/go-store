package store

type StoreName string

const (
	Local StoreName = "local"
	Amazon StoreName = "amazon"
	Google StoreName = "goolge"
	Azure StoreName = "azure"
	Cloudflare StoreName = "cloudflare"
	Dropbox StoreName = "dropbox"
	Qiniu StoreName = "qiniu"
	Aliyun StoreName = "aliyun"
	Minio StoreName = "minio"
)