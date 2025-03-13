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
          name = "seshCore";
          src = ./src;
          vendorHash = "sha256-i+4jsy3utwO8DlngdOmUEpX3Azi1ydHsDhqnwbBhk4c=";
        };
      in {
        default = pkgs.stdenvNoCC.mkDerivation {
          inherit (common) meta version;
          name = "sesh";
          src = ./.;
          dontConfigure = true;
          dontBuild = true;
          propagatedBuildInputs = [ seshCore pkgs.fzf pkgs.tmux ];
          installPhase = ''
            mkdir -p $out/bin $out/lib
            cp $src/hooks/shellHook $out/lib
          '';
          shellHook = ''
            source $src/hooks/shellHook
          '';
        };

      });
  };
}
