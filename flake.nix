{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=24.05";
  };

  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
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
      formatter.x86_64-linux = pkgs.nixpkgs-fmt;
      devShells.x86_64-linux.default = pkgs.mkShell {
        packages = with pkgs; [
          go
        ] ++ [ pulumiBundle ];
      };
    };
}
