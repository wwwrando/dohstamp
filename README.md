# dohstamp

This program: 

- Given an DoH sdns records, outputs an sdns record with an up to certificate. 

In addition, the program:

- Checks that the certificate is valid for the advertised domains
- Displays additional certificate data

# Usage 

To update an old sdns record, say "sdns://AgMAAAAAAAAABzkuOS45Ljkg4-E9rvH9MBLbgLOwArXSp_JKfIu4KzGGlL3K8GHRugISZG5zOS5xdWFkOS5uZXQ6NDQzCi9kbnMtcXVlcnk", 

``./dohstamp -sdns "sdns://AgMAAAAAAAAABzkuOS45Ljkg4-E9rvH9MBLbgLOwArXSp_JKfIu4KzGGlL3K8GHRugISZG5zOS5xdWFkOS5uZXQ6NDQzCi9kbnMtcXVlcnk"``

# Compiling

Either 

``nix build``

with the result in result/bin/dohstamp. Or 

``go build``

with the result in ./main in the current directory. 
