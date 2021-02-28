package docker

import (
	"github.com/docker/docker/pkg/archive"
	"github.com/pkg/errors"
	"io"
)

//func createTarball(src string) (io.Reader, error) {
//	buff := bytes.NewBuffer(nil)
//	tw := tar.NewWriter(buff)
//
//	if err := os.Chdir(src); err != nil {
//		return nil, errors.Wrap(err, "could not change actual dir")
//	}
//
//	// walk through every file in the folder
//	if err := filepath.Walk(".", func(file string, fi os.FileInfo, err error) error {
//		// generate tar header
//		header, err := tar.FileInfoHeader(fi, file)
//		if err != nil {
//			return err
//		}
//
//		// must provide real name
//		// (see https://golang.org/src/archive/tar/common.go?#L626)
//		header.Name = filepath.ToSlash(file)
//
//		// write header
//		if err := tw.WriteHeader(header); err != nil {
//			return err
//		}
//		// if not a dir, write file content
//		if !fi.IsDir() {
//			data, err := os.Open(file)
//			if err != nil {
//				return err
//			}
//			if _, err := io.Copy(tw, data); err != nil {
//				return err
//			}
//		}
//		return nil
//	}); err != nil {
//		return nil, err
//	}
//
//	// produce tar
//	if err := tw.Close(); err != nil {
//		return nil, err
//	}
//
//	fileToWrite, err := os.OpenFile("/opt/debug/compress.tar", os.O_CREATE|os.O_RDWR, 0777)
//
//	if err != nil {
//		panic(err)
//	}
//	if _, err := io.Copy(fileToWrite, buff); err != nil {
//		panic(err)
//	}
//
//	defer fileToWrite.Close()
//
//	return buff, nil
//}

func createTarball(src string) (io.Reader, error) {
	tar, err := archive.TarWithOptions(src, &archive.TarOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not create tar")
	}

	return tar, nil
}
