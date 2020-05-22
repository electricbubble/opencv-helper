package opencv_helper

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var storeDirectory string

func StoreDirectory(pathname string) error {
	if fInfo, err := os.Stat(pathname); err != nil {
		return err
	} else if !fInfo.IsDir() {
		return fmt.Errorf("must be directory '%s'", pathname)
	}
	storeDirectory = pathname
	return nil
}

func checkStoreDirectory() error {
	if storeDirectory == "" {
		return errors.New(`call 'StoreDirectory("/path/dir")' first`)
	}
	return nil
}

var iterationNumber uint32 = 86400

// GenFilename Generate filename in the format `UnixNano() + iterationNumber + Int31n(999).png`
func GenFilename() string {
	unixNano := time.Now().UnixNano()
	rand.Seed(unixNano)
	atomic.CompareAndSwapUint32(&iterationNumber, 86400, 0)
	atomic.AddUint32(&iterationNumber, 1)
	return strconv.FormatInt(unixNano, 10) +
		strconv.FormatUint(uint64(atomic.LoadUint32(&iterationNumber)), 10) +
		strconv.FormatInt(int64(rand.Int31n(999)), 10) + ".png"
}
