{
  buildGoModule,
  fetchFromGitHub,
}:

buildGoModule rec {
  pname = "hizla-waitlist-backend";
  version = "0.0.6";

  src = fetchFromGitHub {
    owner = "hizla";
    repo = "waitlist-backend";
    rev = "refs/tags/v${version}";
    hash = "sha256-3y8eOlrkVnKPVlXwnBHGn7ICNt7T0iGff0Qih7+CtVw=";
  };

  vendorHash = "sha256-g0mkA5S0TFVg1CYXD7NdybER7BYorXXPLpKtxoMJs9A=";

  postInstall = ''
    mv "$out/bin/backend" "$out/bin/${pname}"
  '';
}
