package integrationtest

// import (
// 	"cloud.google.com/go/storage"
// 	"context"
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"crypto/x509"
// 	"encoding/pem"
// 	"fmt"
// 	"github.com/nrf110/integration-test/pkg/gcs"
// 	"github.com/stretchr/testify/assert"
// 	"io"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"testing"
// 	"time"
// )

// func generatePEM() []byte {
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil {
// 		fmt.Println("Failed to generate private key:", err)
// 		os.Exit(1)
// 	}

// 	// Encode the private key in PEM format
// 	privateKeyPEM := &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
// 	}

// 	return pem.EncodeToMemory(privateKeyPEM)
// }

// func TestGCSDependency(t *testing.T) {
// 	init := func(t *testing.T) (context.Context, *storage.BucketHandle) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 		t.Cleanup(cancel)

// 		projectID := "test-project"

// 		dep := gcs.NewDependency()
// 		err := dep.Start(ctx)
// 		assert.NoError(t, err)
// 		t.Cleanup(func() {
// 			assert.NoError(t, dep.Stop(ctx))
// 		})

// 		client := dep.Client().(*storage.Client)

// 		bucket := client.Bucket("test")
// 		err = bucket.Create(ctx, projectID, &storage.BucketAttrs{})
// 		assert.NoError(t, err)

// 		return ctx, bucket
// 	}

// 	t.Run("can upload and download a file", func(t *testing.T) {
// 		ctx, bucket := init(t)

// 		object := bucket.Object("test.txt")
// 		writer := object.NewWriter(ctx)
// 		_, err := io.WriteString(writer, "Hello World")
// 		assert.NoError(t, err)
// 		err = writer.Close()
// 		assert.NoError(t, err)

// 		reader, err := object.NewReader(ctx)
// 		assert.NoError(t, err)
// 		t.Cleanup(func() {
// 			assert.NoError(t, reader.Close())
// 		})
// 		bytes, err := io.ReadAll(reader)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "Hello World", string(bytes))
// 	})

// 	t.Run("can upload and download from a signed url", func(t *testing.T) {
// 		_, bucket := init(t)

// 		uploadUrl, err := bucket.SignedURL("test.txt", &storage.SignedURLOptions{
// 			GoogleAccessID: "test",
// 			Method:         "PUT",
// 			PrivateKey:     generatePEM(),
// 			Expires:        time.Now().Add(5 * time.Minute),
// 			Insecure:       true,
// 			Scheme:         storage.SigningSchemeV4,
// 		})
// 		assert.NoError(t, err)

// 		req, err := http.NewRequest("PUT", uploadUrl, strings.NewReader("Hello World"))
// 		assert.NoError(t, err)

// 		req.Header.Set("Content-Type", "application/octet-stream")
// 		resp, err := http.DefaultClient.Do(req)
// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)

// 		downloadUrl, err := bucket.SignedURL("test.txt", &storage.SignedURLOptions{
// 			GoogleAccessID: "test",
// 			Method:         "GET",
// 			PrivateKey:     generatePEM(),
// 			Expires:        time.Now().Add(5 * time.Minute),
// 			Insecure:       true,
// 			Scheme:         storage.SigningSchemeV4,
// 		})
// 		assert.NoError(t, err)

// 		req, err = http.NewRequest("GET", downloadUrl, nil)
// 		assert.NoError(t, err)
// 		resp, err = http.DefaultClient.Do(req)
// 		assert.NoError(t, err)
// 		assert.Equal(t, resp.StatusCode, http.StatusOK)
// 		bytes, err := io.ReadAll(resp.Body)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "Hello World", string(bytes))
// 	})
// }
