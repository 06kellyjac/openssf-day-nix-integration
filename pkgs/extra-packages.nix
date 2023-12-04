{ pkgs }:
let
  basicpkg = str: pkgs.runCommand "basic-${str}" { } ''
    echo ${str} > $out
  '';
  basicmultipkg = str: pkgs.runCommand "basic-${str}"
    {
      outputs = [ "out" "lib" "extras" ];
    } ''
    echo ${str} > $out
    echo ${str} > $lib
    echo ${str} > $extras
  '';
in
rec {
  # for if a package requires specific args to be replaced
  # somepackage = pkgs.callPackage ./by-name/so/somepackage/package.nix {
  #   something = "hi";
  # };
  hello-static = pkgs.pkgsStatic.hello;
  test-a = basicpkg "a";
  test-b = basicpkg "b";
  test-c = basicmultipkg "c";
  test-d = basicpkg "d";
  test-remote = pkgs.fetchurl {
    name = "hello-remote";
    url = "https://gist.githubusercontent.com/06kellyjac/8f5ebbed9af05263eec5403b046b66f3/raw/0f124b34a2d848ca1ca992c82424eab1a707bef9/hello-world";
    hash = "sha256-Ssq7L3CJPujOEPJo4wkpP5QQ1Y2+U/Hb2X4X9Lo5pqY=";
  };
  test-combined = pkgs.symlinkJoin { name = "test-combined"; paths = [ test-a test-b test-c test-d ]; postBuild = "echo finished all"; };
}
