{
  description = "Project with intelligent skills and agents with claude code";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    skeletons.url = "github:netbrain/skeletons";
  };

  outputs = { self, nixpkgs, flake-utils, skeletons }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            skeletons.packages.${system}.intent-classifier
            # Common development tools
            # Add more tools based on your stack:
            # go, gopls, gotools (for Go)
            # nodejs, typescript (for Node.js/TypeScript)
            # python3, python3Packages.pip (for Python)
            # rustc, cargo (for Rust)
          ];

          shellHook = ''
            # Fix execute permissions on scripts (Nix templates don't preserve +x)
            if [ -d .claude ]; then
              find .claude -name "*.sh" -type f -exec chmod +x {} \; 2>/dev/null
            fi

            echo "🤖 Claude Code environment with intelligent skills & agents"
            echo ""
            echo "Get started:"
            echo "  claude \"help me set up this project with skills and agents\""
            echo ""
          '';
        };
      });
}
