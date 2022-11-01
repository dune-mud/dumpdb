{ pkgs ? (let
  inherit (builtins) fetchTree fromJSON readFile;
  inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
in import (fetchTree nixpkgs.locked) {
  overlays = [ (import "${fetchTree gomod2nix.locked}/overlay.nix") ];
}) }:

let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };

  runCILocally = pkgs.writeScriptBin "ci-local" ''
    echo "Linting..."
    golangci-lint run --verbose .
    echo "Building..."
    nix build
  '';

in pkgs.mkShell {
  packages = [
    goEnv
    runCILocally
    pkgs.gomod2nix
    pkgs.golangci-lint
    pkgs.goreleaser
    pkgs.sqlite
  ];
}
