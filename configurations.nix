{ self }:
{
  vm = import ./nixos/configurations/vm.nix { inherit self; };
}
