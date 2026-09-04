package storage

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

var allowedExt = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

var allowedMIME = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

type SavedFile struct {
	FileName     string
	FileURL      string
	ThumbnailURL string
	MimeType     string
	FileSize     int64
	Width        int
	Height       int
}

type Local struct {
	Root string
}

func NewLocal(root string) (*Local, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &Local{Root: root}, nil
}

func DetectMIME(header *multipart.FileHeader, file multipart.File) (ext string, mime string, err error) {
	ext = strings.ToLower(filepath.Ext(header.Filename))
	mapped, ok := allowedExt[ext]
	if !ok {
		return "", "", fmt.Errorf("不支持的文件类型")
	}
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}
	detected := http.DetectContentType(buf[:n])
	if detected == "application/octet-stream" && ext == ".webp" {
		detected = "image/webp"
	}
	if !allowedMIME[detected] && detected != mapped {
		if !(mapped == "image/jpeg" && (detected == "image/jpeg" || detected == "image/jpg")) {
			if detected != mapped {
				return "", "", fmt.Errorf("文件内容与扩展名不匹配")
			}
		}
	}
	return ext, mapped, nil
}

func (s *Local) Save(file multipart.File, header *multipart.FileHeader, maxBytes int64) (*SavedFile, error) {
	if header.Size > maxBytes {
		return nil, fmt.Errorf("图片不能超过 10MB")
	}
	ext, mime, err := DetectMIME(header, file)
	if err != nil {
		return nil, err
	}
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	day := time.Now().Format("2006-01-02")
	dir := filepath.Join(s.Root, day)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	name := id + ext
	thumbName := id + "_thumb.jpg"
	abs := filepath.Join(dir, name)
	thumbAbs := filepath.Join(dir, thumbName)

	out, err := os.Create(abs)
	if err != nil {
		return nil, err
	}
	written, err := io.Copy(out, file)
	_ = out.Close()
	if err != nil {
		return nil, err
	}

	img, err := imaging.Open(abs, imaging.AutoOrientation(true))
	if err != nil {
		_ = os.Remove(abs)
		return nil, fmt.Errorf("无法解析图片")
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	thumb := imaging.Fit(img, 400, 400, imaging.Lanczos)
	if err := imaging.Save(thumb, thumbAbs, imaging.JPEGQuality(82)); err != nil {
		_ = os.Remove(abs)
		return nil, err
	}

	rel := "/" + filepath.ToSlash(filepath.Join("uploads", day, name))
	thumbRel := "/" + filepath.ToSlash(filepath.Join("uploads", day, thumbName))

	return &SavedFile{
		FileName:     header.Filename,
		FileURL:      rel,
		ThumbnailURL: thumbRel,
		MimeType:     mime,
		FileSize:     written,
		Width:        w,
		Height:       h,
	}, nil
}

func SafeJoin(root, requestPath string) (string, error) {
	clean := strings.TrimPrefix(requestPath, "/uploads/")
	clean = strings.TrimPrefix(clean, "uploads/")
	clean = filepath.Clean(clean)
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("invalid path")
	}
	full := filepath.Join(root, clean)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absFull, absRoot) {
		return "", fmt.Errorf("invalid path")
	}
	return absFull, nil
}
