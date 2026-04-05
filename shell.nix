{
  pkgs ? import <nixpkgs> { },
}:

with pkgs;

mkShellNoCC {
  nativeBuildInputs = [
    gcc
    go
    tailwindcss_3
  ];

  buildInputs = [
    air
  ];

  shellHook = "";
}
