{
  description = "POC for using Zig on Fastly Compute";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    devshell.url = "github:numtide/devshell";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devshell.flakeModule
      ];

      systems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];

      perSystem = { pkgs, system, ... }: {
        devshells.default = {
          commands = [
            {
              name = "serve";
              help = "Start a local dev server";
              command = ''
                # zig build -Doptimize=Debug && \
                fastly compute serve --skip-build --file zig-out/bin/main.wasm
              '';
            }
            {
              name = "publish";
              help = "Compile and publish the WASM executable";
              command = ''
                # zig build -Doptimize=ReleaseSmall && \
                fastly compute pack --wasm-binary zig-out/bin/main.wasm && \
                fastly compute deploy --package pkg/package.tar.gz
              '';
            }
          ];

          packages = with pkgs; [
            fastly
            go_1_22
          ];
        };
      };
    };
}
