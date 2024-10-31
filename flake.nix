{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=24.05";
  };

  outputs = { self, nixpkgs }: let pkgs = nixpkgs.legacyPackages.x86_64-linux; in {
        devShells.x86_64-linux.default = pkgs.mkShell {
          packages = with pkgs; [
            pulumi
            pulumiPackages.pulumi-language-go
            go
          ];
        };
  };
}
