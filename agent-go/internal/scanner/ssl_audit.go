package scanner

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

// ComprehensiveSSLAudit performs a detailed SSL/TLS audit on a target host and port.
func (ns *NetworkScanner) ComprehensiveSSLAudit(host string, port int) (*SSLAudit, error) {
	config := &tls.Config{
		InsecureSkipVerify: true, // We need to connect to analyze the certificate, even if it's invalid.
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for SSL audit: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	audit := &SSLAudit{
		Host: host,
		Port: port,
	}

	// 1. Analyze Certificate Chain
	for _, cert := range state.PeerCertificates {
		c := Certificate{
			Subject:      cert.Subject.String(),
			Issuer:       cert.Issuer.String(),
			ValidFrom:    cert.NotBefore,
			ValidUntil:   cert.NotAfter,
			SelfSigned:   cert.Issuer.String() == cert.Subject.String(),
			SignatureAlg: cert.SignatureAlgorithm.String(),
			SANs:         cert.DNSNames,
		}

		// Check for expiry
		if time.Now().After(cert.NotAfter) {
			audit.Issues = append(audit.Issues, SSLIssue{
				Severity:    "critical",
				Type:        "expired-certificate",
				Description: fmt.Sprintf("Certificate expired on %s", cert.NotAfter.Format(time.RFC1123)),
				Remediation: "Renew the SSL/TLS certificate immediately.",
			})
		} else if cert.NotAfter.Sub(time.Now()) < 30*24*time.Hour {
			audit.Issues = append(audit.Issues, SSLIssue{
				Severity:    "high",
				Type:        "expiring-soon-certificate",
				Description: fmt.Sprintf("Certificate will expire in under 30 days on %s", cert.NotAfter.Format(time.RFC1123)),
				Remediation: "Renew the SSL/TLS certificate soon.",
			})
		}

		// Check for weak key size
		if cert.PublicKeyAlgorithm == x509.RSA {
			if key, ok := cert.PublicKey.(*rsa.PublicKey); ok {
				if key.N.BitLen() < 2048 {
					c.KeySize = key.N.BitLen()
					audit.Issues = append(audit.Issues, SSLIssue{
						Severity:    "high",
						Type:        "weak-key-size",
						Description: fmt.Sprintf("RSA key size of %d bits is considered weak.", key.N.BitLen()),
						Remediation: "Use at least a 2048-bit RSA key or switch to an ECC key.",
					})
				}
			}
		}
		audit.CertificateChain = append(audit.CertificateChain, c)
	}

	// 2. Check Protocol Version
	if state.Version < tls.VersionTLS12 {
		audit.Issues = append(audit.Issues, SSLIssue{
			Severity:    "critical",
			Type:        "weak-protocol",
			Description: fmt.Sprintf("The server supports a weak TLS protocol version: %s.", tls.VersionName(state.Version)),
			Remediation: "Disable TLS 1.0 and 1.1. Enable TLS 1.2 and TLS 1.3.",
		})
	}

	// 3. Check Cipher Suite
	weakCiphers := map[uint16]string{
		tls.TLS_RSA_WITH_RC4_128_SHA:      "RC4",
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA: "3DES",
	}
	if cipherName, isWeak := weakCiphers[state.CipherSuite]; isWeak {
		audit.Issues = append(audit.Issues, SSLIssue{
			Severity:    "high",
			Type:        "weak-cipher",
			Description: fmt.Sprintf("The server is using a weak cipher suite: %s.", cipherName),
			Remediation: "Disable weak cipher suites (like RC4, 3DES) and prioritize modern, secure ciphers (like AES-GCM).",
		})
	}

	// 4. Grade the connection (simplified)
	audit.Grade = "A" // Start with a good grade
	for _, issue := range audit.Issues {
		if issue.Severity == "critical" {
			audit.Grade = "F"
			break
		}
		if issue.Severity == "high" {
			audit.Grade = "B"
		}
	}

	return audit, nil
}
