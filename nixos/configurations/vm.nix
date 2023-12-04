{ self
, nixosSystem ? self.inputs.nixpkgs.lib.nixosSystem
, system ? "x86_64-linux"
}:

nixosSystem {
  inherit system;
  modules =
    [
      ({ ... }: {
        # some stuff to make nix flake check happy
        boot.loader.systemd-boot.enable = true;
        fileSystems."/" = { device = "/dev/disk/by-uuid/made-up"; fsType = "ext4"; };
        system.stateVersion = "23.11";


        # boost the vm's cpu and ram
        virtualisation.vmVariant = {
          virtualisation = {
            cores = 4;
            memorySize = 1024 * 4; # MiB
          };
        };

        # basic users
        users.users.root.initialPassword = "SuperSecret123!!";
        users.users.bob.isNormalUser = true;
        users.users.bob.initialPassword = "password";
        users.users.bob.extraGroups = [ "wheel" ];
        # users.users.jim.isNormalUser = true;
        # users.users.jim.initialPassword = "password";
        # services.clamav.enable = lib.mkOverride 30 false;

        # gnome
        boot.plymouth.enable = true;
        services.xserver.enable = true;
        services.xserver.displayManager.gdm.enable = true;
        services.xserver.desktopManager.gnome.enable = true;
      })
    ];
}
