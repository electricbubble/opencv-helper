package opencv_helper

import (
	"testing"
)

func TestStoreDirectory(t *testing.T) {
	tmpOutput := "checkStoreDirectory()"
	if err := checkStoreDirectory(); err == nil {
		t.Errorf("%s: should be a err", tmpOutput)
	} else {
		t.Logf("[Pass]\t%s: %s", tmpOutput, err)
	}

	tmpOutput = `StoreDirectory("")`
	if err := StoreDirectory(""); err == nil {
		t.Errorf("%s: should be a err", tmpOutput)
	} else {
		t.Logf("[Pass]\t%s: %s", tmpOutput, err)
	}

	tmpOutput = `StoreDirectory("/Users/hero/Documents/Workspace/abc123")`
	if err := StoreDirectory("/Users/hero/Documents/Workspace/abc123"); err == nil {
		t.Errorf("%s: should be a err", tmpOutput)
	} else {
		t.Logf("[Pass]\t%s: %s", tmpOutput, err)
	}

	tmpOutput = `StoreDirectory("/Users/hero/Documents/Workspace/image-helper/helper.go")`
	if err := StoreDirectory("/Users/hero/Documents/Workspace/image-helper/helper.go"); err == nil {
		t.Errorf("%s: should be a err", tmpOutput)
	} else {
		t.Logf("[Pass]\t%s: %s", tmpOutput, err)
	}

	if err := StoreDirectory("/Users/hero/Documents/Workspace/image-helper"); err != nil {
		t.Error(err)
	}

	if err := checkStoreDirectory(); err != nil {
		t.Error(err)
	}

	t.Log(storeDirectory)

}

func Test_genFilename(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(GenFilename())
	}
}
