{ buildGoModule }:

buildGoModule rec {
  pname = "gosh";
  version = "0.0.1";

  src = builtins.filterSource (path: type: type != "directory" || baseNameOf path != ".git") ./.;

  vendorSha256 = "sha256:eGqKXdbCv/9bPuYkfG0gaZ4zn76ddn3oTtZwUIlkRho="; 

  subPackages = [ "." ]; 

  runVend = true;

  buildInputs = [ ];
}


