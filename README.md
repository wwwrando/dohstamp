# dohstamp

This program: 

- Given a DoH sdns records, outputs a sdns record with an up to date certificate. 

In addition, the program:

- Checks that the certificate is valid for the advertised domains
- Displays additional certificate data

# Usage 

To update an old sdns record, e.g "sdns://AgMAAAAAAAAABzkuOS45LjkgKhX11qy258CQGt5Ou8dDsszUiQMrRuFkLwaTaDABJYoSZG5zOS5xdWFkOS5uZXQ6NDQzCi9kbnMtcXVlcnk"

``./dohstamp -sdns "sdns://AgMAAAAAAAAABzkuOS45LjkgKhX11qy258CQGt5Ou8dDsszUiQMrRuFkLwaTaDABJYoSZG5zOS5xdWFkOS5uZXQ6NDQzCi9kbnMtcXVlcnk"``

Additional options, all disabled by default:  

- silent : returns only a working sdns or a one-line error
- show : show the retrieved certificates
- last : hash the certificate provided by the host instead of the certificate that signed the certificate provided by the host (this breaks the recommended usage, but might be more secure)

# Compiling

Either 

``nix build``

with the result in result/bin/dohstamp. Or 

``go build -o dohstamp``

with the result in ./dohstamp in the current directory. 
