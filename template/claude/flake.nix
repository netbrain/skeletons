{
  description = "Project with intelligent skills and agents with claude code";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Build intent-classifier from GitHub
        # Always pulls latest from main branch
        intent-classifier = pkgs.buildGoModule {
          pname = "intent-classifier";
          version = "0.1.0";

          src = pkgs.fetchFromGitHub {
            owner = "netbrain";
            repo = "skeletons";
            rev = "main";
            hash = "sha256-BcxmPrZHhnZjarSnrN51bLiqF8XVxAXWC+5s+oKpLnM=";
          };

          sourceRoot = "source/utils/intent-classifier";

          vendorHash = "sha256-Ks1NEdhqgDRUgN9t3rAv71EmZtxHqUnXP+V+ewRBvoU=";

          buildInputs = [ pkgs.libffi ];

          nativeBuildInputs = [ pkgs.makeWrapper ];

          # Skip tests during build (tests require writable home directory for cache)
          doCheck = false;

          # Set LD_LIBRARY_PATH during build
          preBuild = ''
            export LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath [ pkgs.libffi pkgs.stdenv.cc.cc.lib ]}
          '';

          postInstall = ''
            wrapProgram $out/bin/intent-classifier \
              --prefix LD_LIBRARY_PATH : ${pkgs.lib.makeLibraryPath [ pkgs.libffi pkgs.stdenv.cc.cc.lib ]}
          '';
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            intent-classifier
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
