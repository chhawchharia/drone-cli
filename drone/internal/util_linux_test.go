package internal

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/urfave/cli"
)

func TestNewClientSystemTrust(t *testing.T) {
	if mode := os.Getenv("DRONE_TLS_TEST_MODE"); mode != "" {
		flags := flag.NewFlagSet("tls-test", flag.ContinueOnError)
		flags.String("server", os.Getenv("DRONE_TLS_TEST_SERVER"), "")
		flags.String("token", "test-token", "")
		flags.Bool("skip-verify", mode == "skip-verify", "")
		client, err := NewClient(cli.NewContext(cli.NewApp(), flags, nil))
		if err != nil {
			t.Fatal(err)
		}
		user, err := client.Self()
		if mode == "untrusted" {
			var unknownAuthority x509.UnknownAuthorityError
			if !errors.As(err, &unknownAuthority) {
				t.Fatalf("expected an untrusted certificate error, got %v", err)
			}
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if user.Login != "test-user" {
			t.Fatalf("unexpected user: %#v", user)
		}
		return
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/user" || r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"login":"test-user"}`)
	}))
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	defer server.Close()

	dir := t.TempDir()
	trusted := filepath.Join(dir, "trusted.pem")
	empty := filepath.Join(dir, "empty.pem")
	if err := os.WriteFile(trusted, pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE", Bytes: server.Certificate().Raw,
	}), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(empty, nil, 0600); err != nil {
		t.Fatal(err)
	}
	certDir := filepath.Join(dir, "certs")
	if err := os.Mkdir(certDir, 0700); err != nil {
		t.Fatal(err)
	}

	for _, mode := range []string{"private-ca", "untrusted", "skip-verify"} {
		t.Run(mode, func(t *testing.T) {
			bundle := empty
			if mode == "private-ca" {
				bundle = trusted
			}
			// Each process loads its own system pool; x509 caches roots within a process.
			cmd := exec.Command(os.Args[0], "-test.run=^TestNewClientSystemTrust$", "-test.v")
			cmd.Env = append(os.Environ(),
				"DRONE_TLS_TEST_MODE="+mode,
				"DRONE_TLS_TEST_SERVER="+server.URL,
				"SSL_CERT_FILE="+bundle,
				"SSL_CERT_DIR="+certDir,
			)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%v\n%s", err, output)
			}
		})
	}
}
