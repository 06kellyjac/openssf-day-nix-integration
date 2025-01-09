package cosign

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func SignImage(imageName string, digest string) error {
	log.Info().Str("imageName", imageName).Msg("signing")

	// cosign sign --key ./cosign.key --fulcio-url="http://localhost:5555" --rekor-url=http://localhost:3000 ttl.sh/alpine/alpine:1m@sha256:d695c3de6fcd8cfe3a6222b0358425d40adfd129a8a47c3416faff1a8aece389
	cmd := exec.Command(
		"/nix/store/zh3rkq12cwz44g77b95l1vd6fir360i6-cosign-2.2.1/bin/cosign",
		"sign",
		"--key",
		"/home/jk/projects/personal/openssf-day-nix-integration/cosign.key",
		"--fulcio-url=http://localhost:5555",
		"--rekor-url=http://localhost:3000",
		fmt.Sprintf("%s@%s", imageName, digest),
		"--yes",
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "COSIGN_PASSWORD=")
	log.Info().Str("imageName", imageName).Msg("signing 2")
	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	log.Info().Str("imageName", imageName).Msg("signing 3")
	if err != nil {
		return fmt.Errorf("failed to sign image: %v", err)
	}
	log.Info().Str("imageName", imageName).Msg("signing 4")
	log.Info().Msg("BELOW")
	log.Info().Msg(b.String())
	log.Info().Msg("ABOVE")
	return nil
}
