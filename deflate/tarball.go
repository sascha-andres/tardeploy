package deflate

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/prometheus/log"

	"path"

	"github.com/pkg/errors"
)

// Tarball extracts the tarball to given destination directory
func Tarball(tarball, destinationDirectory string) error {
	gzipStream, err := os.Open(tarball)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not open tarball %s", tarball))
	}
	defer gzipStream.Close()

	return extractTarGz(destinationDirectory, gzipStream)
}

func extractTarGz(destinationDirectory string, gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return errors.Wrap(err, "Could not create GZIP reader")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path.Join(destinationDirectory, header.Name), 0750); err != nil {
				return errors.Wrap(err, "extractTarGz: Mkdir() failed")
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(path.Join(destinationDirectory, header.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0640)
			if err != nil {
				return errors.Wrap(err, "extractTarGz: OpenFile() failed")
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.Wrap(err, "extractTarGz: Copy() failed")
			}
		default:
			return errors.New(fmt.Sprintf("extractTarGz: unknown type: %s in %s",
				header.Typeflag,
				header.Name))
		}
	}

	return nil
}
