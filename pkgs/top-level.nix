# trimmed down pkgs/top-level/default.nix and pkgs/top-level/stage.nix
{ pkgs
, localSystem ? pkgs.system
, lib ? pkgs.lib
}:
let
  autoCaller = import "${pkgs.path}/pkgs/top-level/by-name-overlay.nix";
  autoCalledPackages = autoCaller ./by-name;

  extraPackages = _final: _prev:
    import ./extra-packages.nix { inherit pkgs; };

  # compose overlays
  toFix = lib.foldl' (lib.flip lib.extends) (_final: { }) [
    # autoCaller produces an overlay that expects self to have callPackage
    # provide that from pkgs
    (final: prev: { inherit (pkgs) callPackage; })
    autoCalledPackages
    extraPackages
  ];

  # # helper to remove packages for unsupported platforms
  # stripUnavailable = lib.filterAttrs
  #   (n: drv: builtins.elem localSystem (drv.meta.platforms or [ ]));

  allPackages = lib.fix toFix;
in
# stripUnavailable allPackages
allPackages
