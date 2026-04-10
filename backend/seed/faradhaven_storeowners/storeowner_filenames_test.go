package faradhaven_storeowners

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestStoreownerPortraitFilesMatchSeeds ensures backend/storeowners contains
// one PNG per vendor, named exactly {Name}.png from AllStoreOwnerSeeds (matches S3 keys).
// Skips if backend/storeowners is missing (e.g. CI without assets).
func TestStoreownerPortraitFilesMatchSeeds(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	seedDir := filepath.Dir(thisFile)
	storeDir := filepath.Join(seedDir, "..", "..", "storeowners")
	storeDir, err := filepath.Abs(storeDir)
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	fi, err := os.Stat(storeDir)
	if err != nil {
		t.Skipf("skip: %s (%v)", storeDir, err)
	}
	if !fi.IsDir() {
		t.Fatalf("not a directory: %s", storeDir)
	}

	for _, v := range AllStoreOwnerSeeds() {
		name := v.Name + ".png"
		path := filepath.Join(storeDir, name)
		st, err := os.Stat(path)
		if err != nil {
			t.Errorf("missing portrait for vendor %q: expected file %q", v.Name, path)
			continue
		}
		if st.IsDir() {
			t.Errorf("expected file for vendor %q, found directory: %s", v.Name, path)
		}
	}
}
