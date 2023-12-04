package nix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func getOnlyKey[V any](singleMap map[string]V) (string, V, error) {
	var (
		key string
		val V
	)
	if len(singleMap) != 1 {
		return key, val, errors.New("unexpected val count")
	}
	for k, val := range singleMap {
		return k, val, nil
	}
	// should be unreachable
	log.Panic().Msg("map was checked to have 1 value but then didn't")
	panic(1)
}

func GetDrv(path string) (string, Drv, error) {
	log.Info().Str("path", path).Msg("getting drv")
	var (
		key string
		drv Drv
	)
	drvs, err := GetDrvs(path, false)
	if err != nil {
		return key, drv, err
	}
	name, d, err := getOnlyKey(drvs)
	if err != nil {
		return key, drv, err
	}
	key = name
	drv = d
	return key, drv, nil
}

// warning: The interpretation of store paths arguments ending in `.drv` recently changed. If this command is now failing try again with '/nix/store/path.drv^*'
func GetDrvs(path string, recursive bool) (Drvs, error) {
	log.Info().Str("path", path).Msg("getting drv")
	// "--recursive",
	cmdSlice := []string{"nix", "derivation", "show", path}
	// log.Info().Strs("cmdslice", cmdSlice).Msg("cmdslice drv")
	if recursive {
		cmdSlice = append(cmdSlice, "--recursive")
	}
	// first item as cmd and just following items as args
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get derivation output: %v", err)
	}
	output := b.Bytes()

	var d Drvs
	if err := json.Unmarshal([]byte(output), &d); err != nil {
		return nil, fmt.Errorf("\n%v: %v\n", string(output), err)
	}
	return d, nil
}

func SRITo32(sri string) (string, error) {
	log.Info().Str("sri", sri).Msg("converting sri")
	cmd := exec.Command("nix", "hash", "to-base32", sri)

	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get derivation output: %v", err)
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func selectDrvPathInfo(pathInfos []PathInfo) (PathInfo, error) {
	var drvPathInfo []PathInfo
	var pi PathInfo
	if len(pathInfos) < 1 {
		return pi, errors.New("didn't have any path info")
	}
	for _, pi := range pathInfos {
		if strings.HasSuffix(pi.Path, ".drv") {
			drvPathInfo = append(drvPathInfo, pi)
		}
	}
	if len(drvPathInfo) > 1 {
		return pi, errors.New("found more than 1 derivation")
	}
	pi = drvPathInfo[0]
	return pi, nil
}

func GetDrvPathInfo(path string) (PathInfo, error) {
	var pi PathInfo
	pis, err := GetPathInfos(path)
	if err != nil {
		return pi, err
	}

	pathInfo, err := selectDrvPathInfo(pis)
	if err != nil {
		return pi, err
	}
	pi = pathInfo
	return pi, nil
}

func GetPathInfos(path string) ([]PathInfo, error) {
	log.Info().Str("path", path).Msg("getting path info")
	cmd := exec.Command("nix", "path-info", "--json", path)

	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		return []PathInfo{}, fmt.Errorf("failed to get path-info output: %v", err)
	}
	output := b.Bytes()

	var pis PathInfos
	if err := json.Unmarshal([]byte(output), &pis); err != nil {
		return []PathInfo{}, fmt.Errorf("\n%v: %v\n", string(output), err)
	}

	return pis, nil
}
