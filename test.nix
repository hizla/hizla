{
  nixosTest,
  writers,
  writeText,

  system,
  self,
}:

nixosTest {
  name = "hizla";
  nodes.machine =
    { pkgs, ... }:
    {
      users.users = {
        alice = {
          isNormalUser = true;
          description = "Alice Foobar";
          password = "foobar";
          uid = 1000;
        };
      };

      # Automatically login on tty1 as a normal user:
      services.getty.autologinUser = "alice";

      environment.systemPackages = with pkgs; [ self.packages.${system}.hizla ];

      virtualisation.qemu.options = [
        # Increase performance:
        "-smp 8"
      ];

    };

  testScript = ''
    import json

    start_all()
    machine.wait_for_unit("multi-user.target")
    machine.wait_for_unit("getty@tty1.service")

    # To check hizla version:
    hizlaVersion = machine.succeed("hizla version").strip()
    print(hizlaVersion)

    # To check help text:
    print(machine.succeed("hizla"))


    # Check hizla serve behaviour:
    def check_serve(command):
      machine.send_chars(f"HIZLA_API_LISTEN_ADDRESS=:3000 systemd-cat --identifier=hizla {command}\n")
      versionResp = json.loads(machine.wait_until_succeeds("curl -s http://localhost:3000/api/v1", timeout=10))
      if versionResp["version"] != hizlaVersion:
          raise Exception(f"version mismatch: {versionResp["version"]}, want {hizlaVersion}")
      machine.send_key("ctrl-c")
      machine.wait_until_fails("pgrep hizla", timeout=10)


    check_serve("hizla serve")
    check_serve("hizla serve -j ${writeText "serve.json" (builtins.toJSON { address = ":3000"; })}")
    check_serve("hizla serve -t ${writers.writeTOML "serve.toml" { address = ":3000"; }}")
  '';
}
