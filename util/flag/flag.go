package flag

import "flag"

var (
	SafeRoot      string
	WsmpAppKey    string
	WsmpAppSecret string
	WSMPUrl       string
)

func init() {
	flag.StringVar(&SafeRoot, "safe-root", "/Users/sunkechang/Documents/saySafe", "safe-file-root")
	flag.StringVar(&WsmpAppKey, "wsmp-app-key", "bmlaiplatform", "wsmp sdk app key")
	flag.StringVar(&WSMPUrl, "wsmp-url", "http://bdbl-wsmp-qa-01.bdbl.baidu.com:8120", "wsmp url prefix")
	flag.StringVar(&WsmpAppSecret, "wsmp-app-secret", "92ed9a1b1a22db3933aa4889a00e4209ca0498191ed2755b59b4e56ac97779f0", "wsmp sdk app secret")
}
