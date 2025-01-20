{
  description = "Dev shell for pulumi";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
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
        # Note that Go is not provided, because it does not behave correctly inside a Nix shell.
        packages = [ pulumiBundle ];
      };
    };
}
