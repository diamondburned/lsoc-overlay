{ pkgs ? import <nixpkgs> {} }: pkgs.mkShell (import ./deps.nix pkgs)
