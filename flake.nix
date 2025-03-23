{
  description = "Minimal Vim project switcher";

  inputs.utils.url = "github:NewDawn0/nixUtils";

  outputs = { self, utils, ... }: {
    overlays.default = final: prev: {
      sesh = self.packages.${prev.system}.default;
    };
    packages = utils.lib.eachSystem { } (pkgs:
      let
        common = {
          version = "1.0.0";
          meta = {
            description = "Minimal Vim project switcher";
            homepage = "https://github.com/NewDawn0/sesh";
            license = pkgs.lib.licenses.mit;
            maintainers = with pkgs.lib.maintainers; [ NewDawn0 ];
          };
        };
        seshCore = pkgs.buildGoModule {
          inherit (common) meta version;
          name = "sesh-core";
          src = ./src;
          vendorHash = "sha256-i+4jsy3utwO8DlngdOmUEpX3Azi1ydHsDhqnwbBhk4c=";
        };
        seshSource = pkgs.stdenvNoCC.mkDerivation {
          name = "sesh-source";
          inherit (common) meta version;
          src = ./.;
          dontConfigure = true;
          dontBuild = true;
          installPhase =
            "install -D $src/hooks/SOURCE_ME.sh $out/lib/SOURCE_ME";
        };
      in {
        default = pkgs.symlinkJoin {
          name = "sesh";
          inherit (common) meta version;
          paths = with pkgs; [ seshSource seshCore fzf tmux ];
          shellHook = ''
            source $out/lib/SOURCE_ME.sh
          '';
        };
      });
  };
}
