{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = { url = "github:edolstra/flake-compat"; flake = false; };
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    let inherit (flake-utils.lib) eachDefaultSystem mkApp; in
    eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          ourpkgs = self.packages.${system};
        in
        {
          packages = import ./packages.nix { inherit self pkgs; };
          apps = import ./apps.nix { inherit mkApp ourpkgs; };
          devShells.default = pkgs.mkShell {
            name = "devshell";
            env.IS_IN_DEVSHELL = true;
            nativeBuildInputs = with pkgs; [
              cosign
              nix
              sbomnix
              go-task
              k3d
              kubectl
              graphviz
              go_1_21
              rekor-cli
              skopeo
            ];
          };
        }
      ) // {
      nixosConfigurations = import ./configurations.nix { inherit self; };
    };
}
