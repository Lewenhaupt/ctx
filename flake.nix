{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          system = system;
          config.allowUnfree = true;
        };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "ctx";
          version = ".0.0.1";

          src = ./.;

          vendorHash = "sha256-Bxp4bmoqPCygwnHucdWGn9vVwn6PSg3s9UkwTQDtmHc=";

          subPackages = [ "cmd/ctx" ];

          meta = with pkgs.lib; {
            description = "A CLI tool for combining markdown fragments based on tags";
            homepage = "https://github.com/Lewenhaupt/ctx";
            license = licenses.mit; # Update with actual license
            maintainers = [ ];
            mainProgram = "ctx";
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [
            # go
            pkgs.go
            pkgs.gopls
            pkgs.golangci-lint
            pkgs.gofumpt

            # Build tools
            pkgs.git

            # Node.js for commitlint and git hooks
            pkgs.nodejs
            pkgs.nodePackages.npm

            # System libraries

            # Additional tools that might be useful
            pkgs.lld
          ];

        };
      }
    );
}
