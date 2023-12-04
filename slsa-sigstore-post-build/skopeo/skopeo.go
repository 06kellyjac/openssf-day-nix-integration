package skopeo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func CopyToOCI(sri string) (string, error) {
	log.Info().Str("sri", sri).Msg("pushing container")
	// skopeo copy docker-archive:/tmp/busybox.tar.gz docker://ttl.sh/alpine/alpine
	cmd := exec.Command("skopeo", "copy", fmt.Sprintf("docker-archive:%s", path), fmt.Sprintf("docker://ttl.sh/testing/%s", "a"))

	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get derivation output: %v", err)
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}
