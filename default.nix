{ config, lib, pkgs, ... }:

with lib;

let cfg  = config.services.lsoc-overlay;
	deps = import ./deps.nix pkgs;

	package = pkgs.buildGoModule rec {
		name = "lsoc-overlay";
		version = "0.0.2";

		src = ./.;
		vendorSha256 = "14a88jnx60ng0yq01cccr5qj69pxgk4rw5ilsgqi9flfrpci6xd7";

		buildInputs = deps.buildInputs;
		nativeBuildInputs = deps.nativeBuildInputs;
	};

in {
	options.services.lsoc-overlay = {
		enable = mkEnableOption "Enable the lsoc overlay";

		config = mkOption {
			type = types.attrs;
			default = {
				polling_ms   = 1000;
				red_blink_ms = 1200;
				hidden_procs = [];
				num_scanners = 1;
				window = {
					x = 5;
					y = 5;
				};
			};
			description = "Nix attributes to be converted to JSON as config";
		};
	};

	config = mkIf cfg.enable {
		systemd.user.services.lsoc-overlay = 
			let configFile = pkgs.writeText "lsof.json" (builtins.toJSON cfg.config);

			in {
				Unit = {
					Description = "lsoc user daemon";
					After  = [ "graphical-session-pre.target" ];
					PartOf = [ "graphical-session.target" ];
				};
				Install = {
					WantedBy = [ "grpahical-session.target" ];
				};
				Service = {
					Type = "simple";
					ExecStart  = "${package}/bin/lsoc-overlay -c ${configFile}";
					Restart    = "on-failure";
					RestartSec = 5;
				};
			};
	};
}
