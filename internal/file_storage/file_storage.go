package filestorage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"scavenger/internal/models"
)

type FileType string

const (
	Markdowm    FileType = "markdown"
	LabMaterial FileType = "material"
	LabReport   FileType = "report"
)

type FileStorage struct {
	basePath string
	baseURL  string
}

type StoredFile struct {
	Filename string
	Path     string
	URL      string
	Size     int64
}

func New(cfg models.FSConfig) (*FileStorage, error) {
	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		return nil, err
	}

	for _, fileType := range []FileType{Markdowm, LabMaterial, LabReport} {
		path := filepath.Join(cfg.BasePath, string(fileType))
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	}

	return &FileStorage{
		basePath: cfg.BasePath,
		baseURL:  cfg.BaseURL,
	}, nil
}

func (fs *FileStorage) SaveLabFile(ft FileType, lab *models.Lab, file multipart.File, header *multipart.FileHeader) (*StoredFile, error) {
	return fs.SaveFile(ft, strconv.Itoa(lab.DisciplineID), file, header)
}

func (fs *FileStorage) SaveFile(ft FileType, path string, file multipart.File, header *multipart.FileHeader) (*StoredFile, error) {
	filename := fs.sanitizeFilename(header.Filename)

	filePath := fs.getFilePath(ft, path, filename)

	filename = filepath.Base(filePath)

	dirPath := filepath.Dir(filePath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	url := fs.getURL(ft, path, filename)

	return &StoredFile{
		Filename: filename,
		Path: filePath,
		URL: url,
		Size: size,
	}, nil
}

func (fs *FileStorage) sanitizeFilename(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.Trim(base, ext)

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9._-]`)
	safe := reg.ReplaceAllString(nameWithoutExt, "-")

	reg = regexp.MustCompile(`-+`)
	safe = reg.ReplaceAllString(safe , "-")

	safe = strings.Trim(safe, "-")

	if safe == "" {
		return "unnamed"
	}

	return nameWithoutExt + ext
}

func (fs *FileStorage) generateUniqueFileName(filePath string) string {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath
	}

	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	counter := 1
	for {
		newBase := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		newPath := filepath.Join(dir, newBase)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}

		counter++
		if counter > 1000 {
			newBase = fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
			return filepath.Join(dir, newBase)
		}
	}
}

func (fs *FileStorage) getFilePath(ft FileType, path string, fn string) string {
	fp := filepath.Join(fs.basePath, string(ft), path, fn)
	return fs.generateUniqueFileName(fp)
}

func (fs *FileStorage) getURL(ft FileType, path string, fn string) string {
	return fmt.Sprintf(
		"%s/%s/%s/%s",
		strings.TrimSuffix(fs.baseURL, "/"),
		string(ft),
		path,
		fn,
)
}
