{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:gytis-ivaskevicius/flake-utils-plus";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: {
      packages.gosh =
        nixpkgs.legacyPackages.${system}.callPackage ./gosh.nix {};

      defaultPackage = self.packages.${system}.gosh;
    });
}














