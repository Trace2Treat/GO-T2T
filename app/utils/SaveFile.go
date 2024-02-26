package utils

import (
	// "errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveFile menyimpan file yang diunggah ke folder yang ditentukan dan mengembalikan path penyimpanan file.
func SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Mendapatkan path absolut folder penyimpanan
	dir, err := filepath.Abs(filepath.Join("storage", folder))
	if err != nil {
		return "", err
	}

	// Membuat folder jika belum ada
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Membuka file yang diunggah
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Membuat file tujuan di folder penyimpanan
	dst, err := os.Create(filepath.Join(dir, fileHeader.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Menyalin isi file ke file tujuan
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	// Mengembalikan path penyimpanan file
	return filepath.Join(dir, fileHeader.Filename), nil
}
