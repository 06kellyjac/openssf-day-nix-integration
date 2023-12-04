{ lib
, buildGoModule
, fetchFromGitHub
}:

buildGoModule rec {
  pname = "guacone";
  version = "0.3.0";

  src = fetchFromGitHub {
    owner = "guacsec";
    repo = "guac";
    rev = "refs/tags/v${version}";
    hash = "sha256-pZUEfRTm3uz75ZXzSG8/zBhQLBQEgl7ZgQJQD3a4uPQ=";
  };

  vendorHash = "sha256-v7bs5OTSimix7RgI4KNJaD3TNvnGyEAhLvh7t3JQ7z8=";

  subPackages = [ "cmd/guacone" ];

  ldflags = [
    "-s"
    "-w"
    "-X=github.com/guacsec/guac/pkg/version.Version=v${version}"
  ];

  doCheck = false;

  meta = with lib; {
    homepage = "https://guac.sh";
    changelog = "https://github.com/guacsec/guac/releases/tag/v${version}";
    description = "GUAC aggregates software security metadata into a high fidelity graph database";
    longDescription = ''
      Graph for Understanding Artifact Composition (GUAC) aggregates software
      security metadata into a high fidelity graph database—normalizing entity
      identities and mapping standard relationships between them. Querying this
      graph can drive higher-level organizational outcomes such as audit,
      policy, risk management, and even developer assistance.
    '';
    mainProgram = "guacone";
    license = licenses.asl20;
    maintainers = with maintainers; [ jk ];
  };
}
