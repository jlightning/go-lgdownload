package lgdownload

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

type FileWriter struct {
	*os.File
	locker      sync.Mutex
	startTime   time.Time
	byteWritten uint64
}

func (f *FileWriter) WriteAt(p []byte, off int64) (n int, err error) {
	f.locker.Lock()
	elapsed := time.Now().Sub(f.startTime)
	f.byteWritten += uint64(len(p))

	bytePerSec := f.byteWritten
	if uint64(elapsed.Seconds()) > 0 {
		bytePerSec = f.byteWritten / uint64(elapsed.Seconds())
	}

	fmt.Println("speed", humanize.Bytes(bytePerSec), " per sec")
	f.locker.Unlock()
	return f.File.WriteAt(p, off)
}
