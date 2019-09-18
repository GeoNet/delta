package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const globalStrong = `###
### Delivered by puppet
###
# Defines the channel code of the preferred stream used eg by scautopick and
# scrttv. If no component code is given, 'Z' will be used by default.
detecStream = HN

# Defines the location code of the preferred stream used eg by scautopick and
# scrttv.
detecLocid = 20
`

const globalWeak = `###
### Delivered by puppet
###
# Defines the channel code of the preferred stream used eg by scautopick and
# scrttv. If no component code is given, 'Z' will be used by default.
detecStream = EH

# Defines the location code of the preferred stream used eg by scautopick and
# scrttv.
detecLocid = 10
`

func TestGlobal(t *testing.T) {

	globals := map[string]struct {
		global  Global
		content string
	}{
		"global/profile_strong_20": {
			global: Global{
				Location: "20",
				Stream:   "HN",
			},
			content: globalStrong,
		},
		"global/profile_weak_10": {
			global: Global{
				Location: "10",
				Stream:   "EH",
			},
			content: globalWeak,
		},
	}

	for k, g := range globals {
		t.Run("check "+k, func(t *testing.T) {
			d, err := ioutil.TempDir(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := Store(g.global, d); err != nil {
				t.Fatalf("unable to store key output %s: %v", k, err)
			}

			key, err := ioutil.ReadFile(filepath.Join(d, k))
			if err != nil {
				t.Fatalf("unable to read temp key file %s: %v", d, err)
			}

			if string(key) != g.content {
				t.Errorf("contents mismatch %s", k)
			}
		})
	}
}
