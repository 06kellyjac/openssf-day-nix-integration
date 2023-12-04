{ lib
, dockerTools
, pkgs
}:

dockerTools.buildLayeredImage {
  name = "hello";
  config.Cmd = [ "${pkgs.pkgsStatic.hello}/bin/hello" ];
}
