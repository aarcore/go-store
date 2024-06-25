package store

type StoreName string

const (
	Local StoreName = "local"
	Amazon StoreName = "amazon"
	Google StoreName = "goolge"
	Azure StoreName = "azure"
	Cloudflare StoreName = "cloudflare"
	Qiniu StoreName = "qiniu"
	Aliyun StoreName = "aliyun"
)