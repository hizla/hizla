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
          inherit (pkgs) callPackage;
        in
        {
          default = self.packages.${system}.hizla;
          hizla = callPackage ./package.nix { };

          waitlist = callPackage ./extra/waitlist {
            backend = callPackage ./extra/waitlist/backend.nix { };
            frontend = callPackage ./extra/waitlist/frontend.nix;
          };
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
