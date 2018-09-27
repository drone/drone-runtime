package scripts

import (
	"bytes"
	"fmt"
	"strings"
)

// Windows returns the windows powershell script.
func Windows(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			powershellTrace,
			escaped,
			command,
		))
	}
	return fmt.Sprintf(
		powershellScript,
		buf.String(),
	)
}

const powershellScript = `
$ErrorActionPreference = 'Stop';
&cmd /c "mkdir c:\root";
if ($Env:CI_NETRC_MACHINE) {
$netrc=[string]::Format("{0}\_netrc",$Env:HOME);
"machine $Env:CI_NETRC_MACHINE" >> $netrc;
"login $Env:CI_NETRC_USERNAME" >> $netrc;
"password $Env:CI_NETRC_PASSWORD" >> $netrc;
};
[Environment]::SetEnvironmentVariable("CI_NETRC_PASSWORD",$null);
[Environment]::SetEnvironmentVariable("CI_SCRIPT",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_USERNAME",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_PASSWORD",$null);
%s
`

const powershellTrace = `
Write-Output ('+ %s');
& %s; if ($LASTEXITCODE -ne 0) {exit $LASTEXITCODE}
`
