package main

import (
	"bytes"
	"fmt"
	"github.com/jedisct1/go-dnsstamps"
	"encoding/hex"
	"encoding/pem"
	"crypto/x509"
	"crypto/sha256"
	"crypto/tls"
	"flag"
	"time"
)

func GetCertificatesPEM(address string, last bool) (string, error) {
    conn, err := tls.Dial("tcp", address, &tls.Config{
        InsecureSkipVerify: false,
    })
    if err != nil {
        return "", err
    }
    defer conn.Close()
	var b bytes.Buffer

	i := 0;

	for _, cert := range conn.ConnectionState().PeerCertificates {

		/* If last is true skip the last (host) certificate */
		if !last && i == 0 {
			i = i + 1
			continue;
		}

        err := pem.Encode(&b, &pem.Block{
            Type: "CERTIFICATE",
            Bytes: cert.Raw,
        })

        if err != nil {
            return "", err
        }
    }
    return b.String(), nil
}

func main() {

	var sdns string;
	var last bool;
	var show bool;

	flag.StringVar(&sdns, "sdns", "", "DoH server sdns with outdated certificate")
	flag.BoolVar(&last, "last", false, "Hash the last certificate in the chain instead of the certificate that signed the last certificate in the chain")
	flag.BoolVar(&show, "show", false, "Show the retrieved certificates")
	flag.Parse()

	if sdns == "" {
		flag.Usage()
		return
	}

	stamp, _ := dnsstamps.NewServerStampFromString(sdns);
	fmt.Println("Retrieving PEM certificate for " + stamp.ServerAddrStr + "/" + stamp.ProviderName + ":443")

	if stamp.Proto != dnsstamps.StampProtoTypeDoH {
		fmt.Println("This is not a DoH server, refusing to proceed");
		return
	}

	pemData, err := GetCertificatesPEM(stamp.ServerAddrStr, last)
	if err != nil {
		fmt.Println("WARNING: Failed retrieving the certificate using the IP address, trying DNS provider name")
		pemData, err = GetCertificatesPEM(stamp.ProviderName + ":443", last)
		if err != nil {
			fmt.Println("ERROR: Failed retrieving the certificate")
			fmt.Println(err)
			return
		}
	}

	fmt.Println("")

	if show {
		fmt.Println(pemData)
	}

	if err !=  nil {
		fmt.Println("Unable to connect to " + stamp.ServerAddrStr)
		fmt.Println(err);
		return
	}

	block, _ := pem.Decode([]byte(pemData));

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Unable to parse certificate")
		return
	}

	fmt.Println("Certificate issued by " + parsedCert.Issuer.String())
	fmt.Println(parsedCert.IssuingCertificateURL)

	fmt.Println("Valid from " + parsedCert.NotBefore.Format(time.RFC1123))
	fmt.Println("Valid until " + parsedCert.NotAfter.Format(time.RFC1123))


	certPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("WARNING: Unable to fetch system certificate pool")
	}

	chain, err := parsedCert.Verify(x509.VerifyOptions{
		Roots: certPool,
	})

	fmt.Println("-----------------------------------------")

	if !last {

		if len(chain) == 0 {
			fmt.Println("WARNING: Chain of trust from system certificates does not exist")
		}else{
			fmt.Println("Certificate validated using system certificate pool")
		}

	}


	if last {

		var valid bool;
		valid = true;

		for _, w := range parsedCert.DNSNames {
			err = parsedCert.VerifyHostname(w)
			if err != nil {
				fmt.Println("ERROR: Certificated declared to be valid for " + w + " but verification not passed")
				valid = false;
			}
		}

		if valid {
			fmt.Println("Certificate valid for all declared domains:")
			for _, w := range parsedCert.DNSNames {
				fmt.Println(w)
			}
		}

	}

	fmt.Println("-----------------------------------------")

	sha := sha256.Sum256(parsedCert.RawTBSCertificate)

	ints := make([]uint8, 32)

	for i := 0; i < len(sha); i += 1 {
		ints[i] = uint8(sha[i])
	}

	for v,w := range stamp.Hashes {
		stamp.Hashes[v] = ints
		fmt.Println("Previous hash " + hex.EncodeToString(w))
		fmt.Println("New hash " + hex.EncodeToString(ints))
	}

	fmt.Println("Result: " + stamp.String());

}
