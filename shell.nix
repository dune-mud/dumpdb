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
  shellHook = "unset GOROOT; unset GOPATH;";
  packages = [
    goEnv
    runCILocally
    pkgs.gotools
    pkgs.gomod2nix
    pkgs.golangci-lint
    pkgs.goreleaser
    pkgs.sqlite
  ];
}
