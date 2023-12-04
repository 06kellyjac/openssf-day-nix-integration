{ lib
, buildPlatform
, hostPlatform
, pkgsi686Linux
}:
let
  inherit (pkgsi686Linux.minimal-bootstrap) bash tinycc-musl;
  tinycc = tinycc-musl;

  pname = "hello-mini";
  version = "0.0.1";

  src = ./hello.c;
in
bash.runCommand "${pname}-${version}"
{
  inherit pname version;

  nativeBuildInputs = [
    tinycc.compiler
  ];

  meta = with lib; {
    description = "GNU Find Utilities, the basic directory searching utilities of the GNU operating system";
    homepage = "https://www.gnu.org/software/findutils";
    license = licenses.gpl3Plus;
    maintainers = teams.minimal-bootstrap.members;
    platforms = platforms.unix;
  };
} ''
  # Unpack
  cp ${src} hello.c
  tcc -B ${tinycc.libs}/lib -o hello hello.c
  mkdir $out
  mv hello $out
''
