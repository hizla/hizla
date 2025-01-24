{
  projectRootFile = "flake.nix";
  programs = {
    nixfmt.enable = true;
    gofmt.enable = true;
    shellcheck.enable = true;
    shfmt.enable = true;
    actionlint.enable = true;
    mdformat.enable = true;
  };

  settings.global = {
    excludes = [ "LICENSE" ];
  };
}
