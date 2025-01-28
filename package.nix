{
  lib,
  buildGoModule,
}:

buildGoModule rec {
  pname = "hizla";
  version = "0.0.0";

  src = builtins.path {
    name = "hizla-src";
    path = lib.cleanSource ./.;
    filter = path: type: !(type != "directory" && lib.hasSuffix ".nix" path);
  };

  vendorHash = "sha256-g2CPIG+Ws5v4LS5CNv2873rrDCDH4lEq0qd0PiOt6WM=";

  ldflags = lib.attrsets.foldlAttrs (
    ldflags: name: value:
    ldflags ++ [ "-X github.com/hizla/hizla/internal.${name}=${value}" ]
  ) [ "-s -w" ] { Version = "v${version}"; };

  preBuild = ''
    HOME=$(mktemp -d) go generate ./...
  '';
}
