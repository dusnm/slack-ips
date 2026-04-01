package cmd

import (
	"fmt"
	"os"
	"strings"
)

var (
	// Sometimes the simplest of solutions is the best.
	// If the config ever grows any larger, I'll consider
	// a smarter approach like using embed.FS.
	exampleCfg = `
[app]
bind = '0.0.0.0'
port = 3000
domain = 'slack-ips.example.org'
secure = true
behind_proxy = true
signing_secret = 'your_signing_secret'

[slack]
app_id = 'your_app_id'
client_id = 'your_client_id'
client_secret = 'your_client_secret'
signing_secret = 'your_signing_secret'

[db]
path = './db/ips.db'
`
)

func DumpCfg() {
	fmt.Fprintf(os.Stdout, strings.TrimSpace(exampleCfg))
}
