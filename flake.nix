{
  description = "Dev shell for pulumi";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ { self, nixpkgs, flake-parts }:
    flake-parts.lib.mkFlake { inherit inputs; }
      {
        # systems that this flake can be used on
        systems = [ "aarch64-darwin" "x86_64-linux" ];

        perSystem = { config, pkgs, system, ... }:
          let
            pulumiBundle = pkgs.stdenv.mkDerivation {
              name = "pulumi-bundle";
              phases = [ "installPhase" "fixupPhase" ];
              buildInputs = with pkgs; [ pulumi pulumiPackages.pulumi-language-go ];
              installPhase = ''
                mkdir -p $out/bin
                mkdir -p $out/share
                cp ${pkgs.pulumi}/bin/* $out/bin/
                cp -r ${pkgs.pulumi}/share $out/share
                cp ${pkgs.pulumiPackages.pulumi-language-go}/bin/* $out/bin/
              '';
            };
          in
          {
            formatter = pkgs.nixpkgs-fmt;
            devShells.default = pkgs.mkShell {
              packages = with pkgs; [
                go
                pulumiBundle
              ];
            };
          };
      };
}
