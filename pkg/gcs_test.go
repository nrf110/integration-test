package integrationtest_test

import (
	"cloud.google.com/go/storage"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/nrf110/integration-test/pkg/gcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func generatePEM() []byte {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		os.Exit(1)
	}

	// Encode the private key in PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(privateKeyPEM)
}

var _ = Describe("gcs.Dependency", func() {
	var (
		dep    *gcs.Dependency
		client *storage.Client
		bucket *storage.BucketHandle
	)

	BeforeEach(func(ctx SpecContext) {
		projectID := "test-project"

		dep = gcs.NewDependency()
		err := dep.Start(ctx)
		Expect(err).NotTo(HaveOccurred())

		client = dep.Client().(*storage.Client)

		bucket = client.Bucket("test")
		err = bucket.Create(ctx, projectID, &storage.BucketAttrs{})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func(ctx SpecContext) {
		err := dep.Stop(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can upload and download a file", func(ctx SpecContext) {
		object := bucket.Object("test.txt")
		writer := object.NewWriter(ctx)
		_, err := io.WriteString(writer, "Hello World")
		Expect(err).NotTo(HaveOccurred())
		err = writer.Close()
		Expect(err).NotTo(HaveOccurred())

		reader, err := object.NewReader(ctx)
		Expect(err).NotTo(HaveOccurred())
		defer reader.Close()
		bytes, err := io.ReadAll(reader)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bytes)).To(Equal("Hello World"))
	})

	It("can upload and download from a signed url", func(ctx SpecContext) {
		uploadUrl, err := bucket.SignedURL("test.txt", &storage.SignedURLOptions{
			GoogleAccessID: "test",
			Method:         "PUT",
			PrivateKey:     generatePEM(),
			Expires:        time.Now().Add(5 * time.Minute),
			Insecure:       true,
			Scheme:         storage.SigningSchemeV4,
		})
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest("PUT", uploadUrl, strings.NewReader("Hello World"))
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("Content-Type", "application/octet-stream")
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		downloadUrl, err := bucket.SignedURL("test.txt", &storage.SignedURLOptions{
			GoogleAccessID: "test",
			Method:         "GET",
			PrivateKey:     generatePEM(),
			Expires:        time.Now().Add(5 * time.Minute),
			Insecure:       true,
			Scheme:         storage.SigningSchemeV4,
		})
		Expect(err).NotTo(HaveOccurred())

		req, err = http.NewRequest("GET", downloadUrl, nil)
		Expect(err).NotTo(HaveOccurred())
		resp, err = http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		bytes, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bytes)).To(Equal("Hello World"))
	})
})
