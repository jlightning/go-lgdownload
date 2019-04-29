package lgdownload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/sync/errgroup"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Do(ctx context.Context, client *http.Client, url string, file string, n int) error {
	resp, err := ctxhttp.Head(ctx, client, url)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with %d status code", resp.StatusCode)
	}
	fmt.Println(resp.Header)
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return errors.New("server does not support range requests")
	}
	if resp.ContentLength < 0 {
		return errors.New("server sent invalid Content-Length header")
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	fileWriter := &FileWriter{startTime: time.Now(), File: f}

	sectionLen := resp.ContentLength / int64(n)
	wg, ctx := errgroup.WithContext(ctx)
	for off := int64(0); off < resp.ContentLength; off += sectionLen {
		off := off
		lim := off + sectionLen
		if lim >= resp.ContentLength {
			lim = resp.ContentLength
		}
		wg.Go(func() error { return getPart(ctx, client, fileWriter, url, off, lim) })
	}

	doneMonitor := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				fileWriter.WriteMonitorInformation()
			case <-doneMonitor:
				break
			}
		}
	}()
	err = wg.Wait()
	doneMonitor <- struct{}{}
	return err
}

func getPart(ctx context.Context, client *http.Client, w io.WriterAt, url string, off, lim int64) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", off, lim))
	resp, err := ctxhttp.Do(ctx, client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server responded with %d status code, expected %d", resp.StatusCode, http.StatusPartialContent)
	}
	_, err = io.Copy(newSectionWriter(w, off), resp.Body)
	fmt.Println("copy done at offset", off)
	return err
}

func newSectionWriter(w io.WriterAt, off int64) *sectionWriter {
	return &sectionWriter{w, off}
}

type sectionWriter struct {
	w   io.WriterAt
	off int64
}

func (w *sectionWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.WriteAt(p, w.off)
	w.off += int64(n)
	return
}
