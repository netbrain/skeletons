{
  description = "Claude Code project templates with intelligent skills and agent orchestration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    {
      templates = {
        claude = {
          path = ./template/claude;
          description = "Claude Code project with project-init, skill-creator, and agent-creator";
          welcomeText = ''
            # Claude Code Project Template

            You've initialized a Claude Code project!

            ## Next Steps

            1. Enter the development environment:
               nix develop
               # or: direnv allow

            2. Start Claude Code and initialize your project:
               "Let's start a new project"

            The project-init skill will guide you through:
            - Selecting your tech stack
            - Creating project structure
            - Setting up orchestrator and specialist agents
            - Configuring your development workflow

            See README.md for more details.
          '';
        };

        default = self.templates.claude;
      };
    } // flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go development tools
            go
            gopls
            gotools
            go-tools

            # Python development tools
            python3
            python3Packages.pip
            python3Packages.virtualenv

            # Version control
            git
          ];

          shellHook = ''
            echo "🛠️  Development environment loaded"
            echo ""
            echo "Go: $(go version | cut -d' ' -f3)"
            echo "Python: $(python3 --version)"
            echo ""
            echo "📦 Available templates:"
            echo "  nix flake init -t .#claude"
          '';
        };
      });
}
