{
  description = "A session program";

  inputs.utils.url = "github:NewDawn0/nixUtils";

  outputs = { self, utils }: {
    overlays.default = final: prev: {
      pkgs = prev // { sesh = self.packages.${prev.system}.default; };
    };
    packages = utils.lib.eachSystem { } (pkgs:
      let
        version = "0.0.1";
        meta = {
          description = "A session program";
          homepage = "https://github.com/NewDawn0/sesh";
          license = pkgs.lib.licenses.mit;
          maintainers = with pkgs.lib.maintainers; [ NewDawn0 ];
        };
        seshCore = pkgs.buildGoModule {
          inherit meta version;
          name = "seshCore";
          src = ./src;
          vendorHash = "sha256-i+4jsy3utwO8DlngdOmUEpX3Azi1ydHsDhqnwbBhk4c=";
        };
      in {
        default = pkgs.stdenvNoCC.mkDerivation {
          inherit meta version;
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
