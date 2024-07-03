package r2

import (
	"testing"

	"github.com/aarcore/go-store/store"
)

var tempStore store.Store

func TestMain(m *testing.M) {
	tempStore, _ = Init(R2Config{
		AccountId: "",
		AccessKey: "",
		SecretKey: "",
		Bucket: "",
	})

	m.Run()
}

func Test_Put(t *testing.T) {

}

func Test_PutFile(t *testing.T) {

}

func Test_Copy(t *testing.T) {

}

func Test_Url(t *testing.T) {}

func Test_Rename(t *testing.T) {}

func TestGet(t *testing.T) {}

func TestDelete(t *testing.T) {
	tempStore.Delete("shit")
}

func TestExists(t *testing.T) {}

func TestSize(t *testing.T) {}