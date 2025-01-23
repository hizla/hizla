{
  description = "hizla backend";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11-small";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
    }:
    let
      supportedSystems = [
        "aarch64-linux"
        "i686-linux"
        "x86_64-linux"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      treefmtEval = forAllSystems (system: treefmt-nix.lib.evalModule nixpkgsFor.${system} ./treefmt.nix);
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = self.packages.${system}.hizla;
          hizla = pkgs.callPackage ./package.nix { };
        }
      );

      formatter = forAllSystems (system: treefmtEval.${system}.config.build.wrapper);

      checks = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          formatting = treefmtEval.${system}.config.build.check self;

          lint =
            let
              inherit (pkgs) runCommandLocal deadnix statix;
            in
            runCommandLocal "check-lint"
              {
                nativeBuildInputs = [
                  deadnix
                  statix
                ];
              }
              ''
                cd ${./.}

                echo "running deadnix..."
                deadnix --fail

                echo "running statix..."
                statix check .

                touch $out
              '';
        }
      );
    };
}
