package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	a := assert.New(t)

	client := NewClient()
	err := client.Do(context.Background(), &http.Client{}, "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_30mb.mp4", 64)
	a.Nil(err)
}
