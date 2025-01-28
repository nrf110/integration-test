package integrationtest_test

import (
	"cloud.google.com/go/storage"
	"github.com/nrf110/integration-test/pkg/gcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("gcs.Dependency", func() {
	It("can upload and download a file", func(ctx SpecContext) {
		projectID := "test-project"

		dep := gcs.NewDependency()
		err := dep.Start(ctx)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			err = dep.Stop(ctx)
			Expect(err).NotTo(HaveOccurred())
		}()

		client := dep.Client().(*storage.Client)

		bucket := client.Bucket("test")
		err = bucket.Create(ctx, projectID, &storage.BucketAttrs{})
		Expect(err).NotTo(HaveOccurred())

		object := bucket.Object("test.txt")
		writer := object.NewWriter(ctx)
		_, err = io.WriteString(writer, "Hello World")
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
})
