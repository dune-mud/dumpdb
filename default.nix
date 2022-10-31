{ pkgs ? (let
  inherit (builtins) fetchTree fromJSON readFile;
  inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
in import (fetchTree nixpkgs.locked) {
  overlays = [ (import "${fetchTree gomod2nix.locked}/overlay.nix") ];
}), self }:

pkgs.buildGoApplication rec {
  pname = "dumpdb";
  version = "0.0.2";
  pwd = ./.;
  src = ./.;
  ldflags = ''
    -X main.commit=${self.shortRev or "dirty"}
  '';
  modules = ./gomod2nix.toml;
}
