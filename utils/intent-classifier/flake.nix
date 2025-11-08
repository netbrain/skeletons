{
  description = "Intent classifier dev environment";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
    devShells.x86_64-linux.default = pkgs.mkShell {
      buildInputs = [
        pkgs.libffi
        pkgs.stdenv.cc.cc.lib  # Provides libstdc++
      ];
      shellHook = ''
        export LD_LIBRARY_PATH=${pkgs.libffi}/lib:${pkgs.stdenv.cc.cc.lib}/lib:$LD_LIBRARY_PATH
      '';
    };
  };
}
