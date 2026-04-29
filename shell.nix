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

  CGO_ENABLED = 1;
}
