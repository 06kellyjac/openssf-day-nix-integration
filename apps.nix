{ mkApp, ourpkgs }:
{
  vm = mkApp { name = "run-nixos-vm"; drv = ourpkgs.vm; };
}
