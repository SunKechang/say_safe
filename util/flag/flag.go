package flag

import (
	"crypto/rsa"
	"flag"
)

var (
	SafeRoot string
	LogPath  string

	PubKey string
	PriKey *rsa.PrivateKey

	PubPath string
	PriPath string
)

func init() {
	flag.StringVar(&SafeRoot, "safe-root", "./safeFiles", "safe-file-root")
	flag.StringVar(&LogPath, "log-path", "log.txt", "file to write logs")
	flag.StringVar(&PubPath, "pub-path", "./public.key", "rsa public key path")
	flag.StringVar(&PriPath, "pri-path", "./private.key", "rsa private key path")
}
