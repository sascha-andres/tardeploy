// Copyright © 2017 Sascha Andres <sascha.andres@outlook.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deflate

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

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
