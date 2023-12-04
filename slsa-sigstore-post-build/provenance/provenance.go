package provenance

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build/nix"
	"github.com/google/uuid"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
	slsacommon "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	slsav1 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v1"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// http://localhost:5000/29cxziviylsf2rhpijrz9sizl53n9h3s.narinfo
// http://localhost:5000/29cxziviylsf2rhpijrz9sizl53n9h3s.nar
// http://localhost:5000/29cxziviylsf2rhpijrz9sizl53n9h3s.narinfo
// http://localhost:5000/log/gspmd3alks44drwfkx75pm8h381lazmd-hello.tar.gz.drv

const (
	// NixPostHookBuilderID is a default builder ID for this post-build-hook producing SLSA provenance.
	NixPostHookBuilderID = "https://github.com/06kellyjac/openssf-day-nix-integration/slsa-sigstore-post-build@v0"
	BuildTypeGeneric     = "https://nixos.org/buildtype/generic"
	BuildTypeContainer   = "https://nixos.org/buildtype/container"
	BuildTypeFetchURL    = "https://github.com/NixOS/nixpkgs/tree/nixos-unstable/pkgs/build-support/fetchurl"
)

var NixHashRegexp = regexp.MustCompile("\\w{32}")

type NixExternal struct {
	Sandbox bool     `json:"sandbox"`
	Builder string   `json:"builder"`
	Args    []string `json:"args"`
}

type FetchURLInternal struct {
	MirrorFile string `json:"mirrorFile"`
	URL        string `json:"URL"`
}

// TODO: resolve mirror list
// TODO: support multi urls
func fetchURLBuildType(pb nix.PostBuild) (slsav1.ProvenanceBuildDefinition, error) {
	var pbd slsav1.ProvenanceBuildDefinition
	rawURL, ok := pb.Derivation.Env["urls"]
	if !ok {
		return pbd, errors.New("couldn't find the URL")
	}

	mirrorListPath, ok := pb.Derivation.Env["mirrorsFile"]
	if !ok {
		return pbd, errors.New("couldn't find the mirrors file for the mirror list")
	}

	mirrorList, err := os.ReadFile(mirrorListPath)
	if err != nil {
		return pbd, fmt.Errorf("could not read mirror list %q: %v", mirrorListPath, err)
	}
	mirrorListB64 := base64.StdEncoding.EncodeToString(mirrorList)

	URL, err := url.Parse(rawURL)
	if err != nil {
		return pbd, fmt.Errorf("could not parse url %q: %v", rawURL, err)
	}
	splitPath := strings.Split(URL.Path, "/")
	if len(splitPath) < 1 {
		return pbd, fmt.Errorf("path %q didn't split", URL.Path)
	}
	base := splitPath[len(splitPath)-1]
	return slsav1.ProvenanceBuildDefinition{
		BuildType:          BuildTypeFetchURL + "@v0",
		ExternalParameters: genExternal(pb),
		InternalParameters: FetchURLInternal{
			MirrorFile: mirrorListB64,
			URL:        rawURL,
		},
		ResolvedDependencies: []slsav1.ResourceDescriptor{
			{
				Name:   base,
				URI:    rawURL,
				Digest: genDigest(pb.OutputsPathInfo[pb.Derivation.Outputs["out"].Path]),
			},
		},
	}, nil
}

func genExternal(pb nix.PostBuild) NixExternal {
	return NixExternal{
		Sandbox: true,
		Builder: pb.Derivation.Builder,
		Args:    pb.Derivation.Args,
	}
}

func IsContainerBuild(pb nix.PostBuild) bool {
	buildCommand, ok := pb.Derivation.Env["buildCommand"]
	return ok && strings.Contains(buildCommand, " pigz ")
}

func genBuildDefinition(pb nix.PostBuild) slsav1.ProvenanceBuildDefinition {
	buildType := BuildTypeGeneric + "@v0"

	_, ok := pb.Derivation.Env["mirrorsFile"]
	if ok {
		res, err := fetchURLBuildType(pb)
		if err != nil {
			log.Err(err).Msg("failed to process fetchurl, continuing with generic")
		} else {
			return res
		}
	}

	buildCommand, ok := pb.Derivation.Env["buildCommand"]
	if ok && strings.Contains(buildCommand, " pigz ") {
		return slsav1.ProvenanceBuildDefinition{
			ExternalParameters: genExternal(pb),
			BuildType:          BuildTypeContainer + "@v0",
		}
	}

	return slsav1.ProvenanceBuildDefinition{
		ExternalParameters: genExternal(pb),
		BuildType:          buildType,
	}
}

func identifyByproducts(pb nix.PostBuild) []slsav1.ResourceDescriptor {
	// rd := make([]slsav1.ResourceDescriptor, 0)
	// http://localhost:5000/log/gspmd3alks44drwfkx75pm8h381lazmd-hello.tar.gz.drv

	nameNoStore := strings.TrimPrefix(pb.Name, "/nix/store/")
	return []slsav1.ResourceDescriptor{
		{
			Name:      "logs",
			MediaType: "text/plain",
			URI:       fmt.Sprintf("http://localhost:5000/log/%s", nameNoStore),
		},
	}
}

func genDigest(pi nix.PathInfo) slsacommon.DigestSet {
	return slsacommon.DigestSet{
		"SRI": pi.NarHash,
	}
}

func genSubject(pb nix.PostBuild) []intoto.Subject {
	subs := make([]intoto.Subject, 0)
	subs = append(subs, intoto.Subject{
		Name:   pb.Name,
		Digest: genDigest(pb.DerivationPathInfo),
	})
	eg := errgroup.Group{}
	for output, pathInfo := range pb.OutputsPathInfo {
		o := output
		pi := pathInfo
		log.Info().Str("output", o).Any("pi", pi).Msg("hi")

		eg.Go(func() error {
			// TODO: get nar location from narinfo
			nixHash := nixHashRegexp.FindString(pi.Path)
			b32, err := nix.SRITo32(pi.NarHash)
			if err != nil {
				return err
			}
			subs = append(subs, intoto.Subject{
				// Name:   fmt.Sprintf("http://localhost:5000/%s.nar", nixHash),
				// nar/87cky2hh5hzak73c8vsrymb0v7hbw9a8-1fg89d2knmknd2rms9fa4vhx79aw62sf6qd3868gqdzhjqzfyij5.nar
				Name:   fmt.Sprintf("http://localhost:5000/nar/%s-%s.nar", nixHash, b32),
				Digest: genDigest(pi),
			})
			return nil
		})
	}
	eg.Wait()
	return subs
}

// Generate generates an in-toto provenance statement in SLSA v1 format.
func Generate(UUID uuid.UUID, buildStart time.Time, buildFinish time.Time, pb nix.PostBuild) (*intoto.ProvenanceStatementSLSA1, error) {
	return &intoto.ProvenanceStatementSLSA1{
		StatementHeader: intoto.StatementHeader{
			Type:          intoto.StatementInTotoV01,
			PredicateType: slsav1.PredicateSLSAProvenance,
			Subject:       genSubject(pb),
		},
		Predicate: slsav1.ProvenancePredicate{
			BuildDefinition: genBuildDefinition(pb),
			RunDetails: slsav1.ProvenanceRunDetails{
				Builder: slsav1.Builder{
					ID: NixPostHookBuilderID,
				},
				BuildMetadata: slsav1.BuildMetadata{
					InvocationID: UUID.String(),
					FinishedOn:   &buildFinish,
				},
				Byproducts: identifyByproducts(pb),
			},
			// Invocation:  invocation,
			// BuildConfig: buildConfig,
			// Materials:   materials,
			// Metadata:    metadata,
		},
	}, nil
}
