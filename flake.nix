{
  description = "A Nix-flake-based Go 1.17 development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.flake-compat.url = "github:edolstra/flake-compat";
  inputs.flake-compat.flake = false;

  # Closes commit in foundry.nix to forge 94ae897 used in CI.
  inputs.foundry.url = "github:shazow/foundry.nix/e541c43c012d16c6a6f2a8a7a7145dc0f76f004f";


  outputs = { self, flake-utils, nixpkgs, foundry, ... }:
    let
      goVersion = 21; # Change this to update the whole stack
      overlays = [
        (final: prev: {
          go = prev."go_1_${toString goVersion}";
          # Overlaying nodejs here to ensure nodePackages use the desired
          # version of nodejs.
          nodejs = prev.nodejs-20_x;
          pnpm = prev.nodePackages.pnpm;
          yarn = prev.nodePackages.yarn;
        })
        foundry.overlay
      ];
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit overlays system;
          config = {
            permittedInsecurePackages = [ "nodejs-16.20.2" ];
          };
        };
      in
      {
        devShells.default = pkgs.mkShell {
          COMPOSE_DOCKER_CLI_BUILD=1;
          DOCKER_BUILDKIT=1;
          packages = with pkgs; [
            go
            # goimports, godoc, etc.
            gotools
            # https://github.com/golangci/golangci-lint
            golangci-lint

            # Node
            nodejs
            pnpm
            yarn # `pnpm build` fails without this

            # Foundry, and tools like the anvil dev node
            foundry-bin
            solc

            # Docker
            docker-compose # provides the `docker-compose` command

            # Python
            # (python3.withPackages (ps: with ps; [ ]))
            (python3.withPackages (python-pkgs: [
              python-pkgs.requests
              python-pkgs.pyjwt
            ]))
            jq

            # geth node
            go-ethereum
          ] ++ lib.optionals stdenv.isDarwin [
            darwin.libobjc
            darwin.IOKit
            darwin.apple_sdk.frameworks.CoreFoundation
          ];
        };
      });
}
