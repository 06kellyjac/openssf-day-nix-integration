package nix

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type PostBuild struct {
	Name               string
	Derivation         Drv
	DerivationPathInfo PathInfo
	OutputsPathInfo    map[string]PathInfo
}

func CollectOutputsPathInfos(paths []string) (map[string]PathInfo, error) {
	eg := errgroup.Group{}
	output := make(map[string]PathInfo, 3)
	for _, path := range paths {
		p := path
		eg.Go(func() (err error) {
			pis, err := GetPathInfos(p)
			if err != nil {
				return fmt.Errorf("failed to fetch path info for %q: %v", p, err)
			}
			if len(pis) != 1 {
				return fmt.Errorf("got more than 1 path info for what is meant to be an output: %q", len(pis))
			}
			var pi PathInfo
			for _, pathinfo := range pis {
				pi = pathinfo
			}
			output[p] = pi
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return output, err
	}
	return output, nil
}

func NewPostBuild(derivation string, outputs []string) (*PostBuild, error) {
	name, drv, err := GetDrv(derivation)
	if err != nil {
		return nil, err
	}

	drvPI, err := GetDrvPathInfo(name)
	if err != nil {
		return nil, err
	}

	outputPathInfo, err := CollectOutputsPathInfos(outputs)
	if err != nil {
		return nil, err
	}

	return &PostBuild{
		Name:               name,
		Derivation:         drv,
		DerivationPathInfo: drvPI,
		OutputsPathInfo:    outputPathInfo,
	}, nil
}
