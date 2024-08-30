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
