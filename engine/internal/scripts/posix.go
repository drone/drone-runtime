package scripts

import (
	"bytes"
	"fmt"
	"strings"
)

// Posix returns the posix shell script
func Posix(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			shellTrace,
			escaped,
			command,
		))
	}
	return fmt.Sprintf(
		shellScript,
		buf.String(),
	)
}

const shellScript = `
if [ -n "$CI_NETRC_MACHINE" ]; then
cat <<EOF > $HOME/.netrc
machine $CI_NETRC_MACHINE
login $CI_NETRC_USERNAME
password $CI_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
fi
unset CI_NETRC_USERNAME
unset CI_NETRC_PASSWORD
unset CI_SCRIPT
unset DRONE_NETRC_USERNAME
unset DRONE_NETRC_PASSWORD

set -e

%s
`

const shellTrace = `
echo + %s
%s
`
