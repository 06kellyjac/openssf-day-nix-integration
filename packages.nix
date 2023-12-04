{ self, pkgs }:
let
  ourpkgs = import ./pkgs/top-level.nix { inherit pkgs; };
  nixpkgs = pkgs;
in
# our packages + flake specific extras
ourpkgs // {
  vm = self.nixosConfigurations.vm.config.system.build.vm;
}
