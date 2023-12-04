package nix

type InputDrv struct {
	Outputs []string `json:"outputs"`
}

type InputSrcs []string

type Output struct {
	Hash     string `json:"hash"`
	HashAlgo string `json:"hashAlgo"`
	Path     string `json:"path"`
}

type Drv struct {
	Args      []string            `json:"args"`
	Builder   string              `json:"builder"`
	Env       map[string]string   `json:"env"`
	InputDrvs map[string]InputDrv `json:"inputDrvs"`
	InputSrcs InputSrcs           `json:"inputSrcs"`
	Name      string              `json:"name"`
	Outputs   map[string]Output   `json:"outputs"`
	System    string              `json:"system"`
}

type Drvs map[string]Drv

type PathInfo struct {
	Deriver          string   `json:"deriver"`
	NarHash          string   `json:"narHash"`
	NarSize          int      `json:"narSize"`
	Path             string   `json:"path"`
	References       []string `json:"references"`
	RegistrationTime int      `json:"registrationTime"`
	Ultimate         bool     `json:"ultimate"`
	Signatures       []string `json:"signatures"`
	Valid            bool     `json:"valid"`
}

type PathInfos []PathInfo

type PathInfoDrvOutMap map[string]map[string]PathInfos
