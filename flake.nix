{
  description = "A Go command-line tool that converts Brakeman security scan results to GitLab Code Quality format";

  nixConfig = {
    extra-substituters = [
      "https://omochice.cachix.org"
    ];
    extra-trusted-public-keys = [
      "omochice.cachix.org-1:d+cdfbGVPgtxxdGSkGf3hhaCdfziMtZ6FSHUWxwUTo8="
    ];
  };

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixpkgs-unstable";
    nur-packages = {
      url = "github:Omochice/nur-packages";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      nur-packages,
      flake-utils,
      treefmt-nix,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            nur-packages.overlays.default
          ];
        };
        treefmt = treefmt-nix.lib.evalModule pkgs (
          { ... }:
          {
            settings.global.excludes = [
              "CHANGELOG.md"
            ];
            settings.formatter = {
              # keep-sorted start block=yes
              rumdl = {
                command = "${pkgs.lib.getExe pkgs.rumdl}";
                options = [
                  "fmt"
                  "--config"
                  (toString ./rumdl.toml)
                ];
                includes = [ "*.md" ];
              };
              # keep-sorted end
            };
            programs = {
              # keep-sorted start block=yes
              formatjson5 = {
                enable = true;
                indent = 2;
              };
              gofmt.enable = true;
              goimports.enable = true;
              keep-sorted.enable = true;
              nixfmt.enable = true;
              taplo = {
                enable = true;
              };
              yamlfmt = {
                enable = true;
                settings = {
                  formatter = {
                    type = "basic";
                    retain_line_breaks_single = true;
                  };
                };
              };
              # keep-sorted end
            };
          }
        );
        version = pkgs.lib.pipe ./.github/release-please-manifest.json [
          builtins.readFile
          builtins.fromJSON
          (builtins.getAttr ".")
        ];
        runAs =
          name: runtimeInputs: text:
          let
            program = pkgs.writeShellApplication {
              inherit name runtimeInputs text;
            };
          in
          {
            type = "app";
            program = "${program}/bin/${name}";
          };
        devPackages = rec {
          actions = [
            pkgs.actionlint
            pkgs.ghalint
            pkgs.zizmor
          ];
          default = actions ++ [
            pkgs.go_1_26
            pkgs.goreleaser
            treefmt.config.build.wrapper
          ];
        };
      in
      {
        # keep-sorted start block=yes
        apps = {
          check-actions = pkgs.lib.pipe ''
            actionlint
            ghalint run
            zizmor .github/workflows
          '' [ (runAs "check-actions" devPackages.actions) ];
        };
        checks = {
          formatting = treefmt.config.build.check self;
        };
        devShells = pkgs.lib.pipe devPackages [
          (pkgs.lib.attrsets.mapAttrs (name: buildInputs: pkgs.mkShell { inherit buildInputs; }))
        ];
        formatter = treefmt.config.build.wrapper;
        packages = {
          default = pkgs.buildGo126Module {
            #keep-sorted start block=yes
            pname = "brakeman-to-codequality";
            version = version;
            src = ./.;
            vendorHash = "sha256-x0Xzh7SYDE4mSwTl2XFHeZ+CqB6hzzeJcwNXYMEo6q0=";
            env.CGO_ENABLED = 0;
            ldflags = [
              "-X main.version=${version}"
            ];
            meta.description = "A Go command-line tool that converts Brakeman security scan results to GitLab Code Quality format";
            meta.homepage = "https://github.com/Omochice/brakeman-to-codequality";
            meta.license = pkgs.lib.licenses.zlib;
            #keep-sorted end
          };
        };
        # keep-sorted end
      }
    );
}
