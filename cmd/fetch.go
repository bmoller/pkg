/*
 * fetch.go
 *
 * Copyright (c) 2024 Brandon Moller
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/spf13/cobra"

	"github.com/bmoller/pkg/aur"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch package",
	Short: "Fetch a package snapshot from the AUR",
	Long: `The fetch command retrieves a snapshot of the requested package from the AUR. The
archive is saved to a temporary directory and extracted in the current location,
unless otherwise specified.`,
	Args: cobra.ExactArgs(1),
	Run:  fetch,
}

func fetch(cmd *cobra.Command, args []string) {
	filepath, err := aur.DownloadSnapshot(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("failed to open downloaded archive: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		fmt.Printf("failed to decompress tarball: %s\n", err)
		os.Exit(1)
	}
	defer gzReader.Close()
	tarReader := tar.NewReader(gzReader)

	// iterate over tarball contents and extract the important bits
	// we really only care about files, directories, and symlinks
	for h, err := tarReader.Next(); err != io.EOF; h, err = tarReader.Next() {
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(h.Name, fs.FileMode(h.Mode)); err != nil {
				fmt.Printf("failed to create directory '%s': %s\n", h.Name, err)
				os.Exit(1)
			}
		case tar.TypeReg:
			if t, err := os.OpenFile(h.Name, os.O_CREATE|os.O_RDWR, fs.FileMode(h.Mode)); err != nil {
				fmt.Printf("failed to open output file '%s': %s\n", h.Name, err)
				os.Exit(1)
			} else if _, err := t.ReadFrom(tarReader); err != nil && !errors.Is(err, io.EOF) {
				fmt.Printf("failed to read data for output file '%s': %s\n", h.Name, err)
				os.Exit(1)
			} else if err := t.Close(); err != nil {
				fmt.Printf("failed to close output file '%s': %s\n", h.Name, err)
				os.Exit(1)
			}
		case tar.TypeSymlink:
			if err := os.Symlink(h.Linkname, h.Name); err != nil {
				fmt.Printf("failed to create symlink '%s': %s\n", h.Name, err)
				os.Exit(1)
			}
		}
	}
}
