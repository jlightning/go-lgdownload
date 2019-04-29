package lgdownload

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 545e11bab5aa685bb6f8f5092be87199
func TestDo(t *testing.T) {
	a := assert.New(t)

	url := "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_30mb.mp4"
	fileUrl := "./test/big_buck_bunny_720p_30mb.mp4"

	//checksum := "545e11bab5aa685bb6f8f5092be87199"

	client := NewClient()
	err := client.Do(context.Background(), &http.Client{}, url, fileUrl, 8)
	a.Nil(err)

	bytes, err := ioutil.ReadFile(fileUrl)
	h := md5.New()
	h.Write(bytes)

	checkSum1 := fmt.Sprintf("%x", h.Sum(nil))

	client = NewClient()
	err = client.Do(context.Background(), &http.Client{}, url, fileUrl, 1)
	a.Nil(err)

	bytes, err = ioutil.ReadFile(fileUrl)
	h = md5.New()
	h.Write(bytes)

	checkSum2 := fmt.Sprintf("%x", h.Sum(nil))

	a.Equal(checkSum1, checkSum2)
}
