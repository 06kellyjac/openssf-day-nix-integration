package skopeo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build/provenance"
	"github.com/rs/zerolog/log"
)

const base = "ttl.sh/testing/post-hook"

func CopyToOCI(path string) (string, string, error) {
	log.Info().Str("path", path).Msg("pushing container")
	nixHash := provenance.NixHashRegexp.FindString(path)
	imageName := fmt.Sprintf("%s:nix-%s", base, nixHash)
	withTransport := fmt.Sprintf("docker://%s", imageName)
	// skopeo copy docker-archive:/tmp/busybox.tar.gz docker://ttl.sh/alpine/alpine
	cmd := exec.Command("/nix/store/yvi7p645dg4gg05g941rn9xi3kh6kik7-skopeo-1.13.3/bin/skopeo", "copy", fmt.Sprintf("docker-archive:%s", path), withTransport)
	err := cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to push image: %v", err)
	}

	// skopeo inspect docker://$IMAGE_NAME --format "{{ .Digest }}"
	digestCmd := exec.Command("/nix/store/yvi7p645dg4gg05g941rn9xi3kh6kik7-skopeo-1.13.3/bin/skopeo", "inspect", withTransport, "--format", "{{ .Digest }}")
	var b bytes.Buffer
	digestCmd.Stdout = &b
	err = digestCmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to get digest: %v", err)
	}
	return imageName, strings.TrimSuffix(b.String(), "\n"), nil
}
