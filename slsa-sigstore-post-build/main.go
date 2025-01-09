package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build/nix"
	"github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build/provenance"
	cosign "github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build/sigstore"
	"github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build/skopeo"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TODO: catch panic and delete store paths
// the hook failing stops the nix builds so best to avoid failure or properly clean up

func cmd(UUID uuid.UUID, dryrun bool) error {
	buildFinished := time.Now()
	log.Info().Msg("processing")
	drvPath := os.Getenv("DRV_PATH")
	outPathsRaw := os.Getenv("OUT_PATHS")
	outPaths := strings.Split(outPathsRaw, " ")

	log.Trace().Str("drvpath", drvPath).Strs("outpaths", outPaths).Msg("fetching data")
	pb, err := nix.NewPostBuild(drvPath, outPaths)
	if err != nil {
		return err
	}

	icb := provenance.IsContainerBuild(*pb)
	if icb {
		if dryrun {
			log.Trace().Msg("not pushing")
		} else {
			log.Info().Msg("pushing")
			imageName, digest, err := skopeo.CopyToOCI((*pb).Derivation.Outputs["out"].Path)
			if err != nil {
				return err
			}
			err = cosign.SignImage(imageName, digest)
			if err != nil {
				return err
			}
		}
	}

	log.Trace().Str("drvpath", drvPath).Strs("outpaths", outPaths).Msg("transforming data")
	provenance, err := provenance.Generate(UUID, buildFinished, buildFinished, *pb)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Trace().Str("drvpath", drvPath).Strs("outpaths", outPaths).Msg("provenance generated")
	b, err := json.Marshal(provenance)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	dryRun := false
	// user := os.Getenv("USER")
	UUID := uuid.New()
	logfile, _ := os.OpenFile(
		"/tmp/post-build.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if !dryRun {
		multi := zerolog.MultiLevelWriter(os.Stdout, logfile)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Logger = log.With().Str("uuid", UUID.String()).Logger()
	err := cmd(UUID, dryRun)
	if err != nil {
		log.Err(err).Msg("hook failed")
	}

}
