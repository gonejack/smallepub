package smallepub

import (
	"archive/zip"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SmallEpub struct {
	Options
}

func (h *SmallEpub) Run() (err error) {
	if h.About {
		fmt.Println("Visit https://github.com/gonejack/smallepub")
		return
	}
	if len(h.EPUB) == 0 {
		return errors.New("no .epub file given")
	}
	h.run()
	return
}

func (h *SmallEpub) run() {
	for _, epub := range h.EPUB {
		if strings.HasSuffix(epub, ".small.epub") {
			continue
		}
		out := strings.TrimSuffix(epub, filepath.Ext(epub)) + ".small.epub"
		log.Printf("processing %s => %s", epub, out)
		_, exx := os.Stat(out)
		if exx == nil {
			log.Printf("output file %s already exist", out)
			continue
		}
		err := h.do(epub, out)
		if err != nil {
			log.Printf("process %s failed: %s", epub, err)
		}
	}
}

func (h *SmallEpub) do(input string, output string) (err error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "smallepub_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	err = unzip(input, tempDir)
	if err != nil {
		return fmt.Errorf("cannot unzip file %s: %w", input, err)
	}
	err = processImages(tempDir, h.Quality)
	if err != nil {
		return fmt.Errorf("cannot process picture files: %w", err)
	}
	err = zipDir(tempDir, output)
	if err != nil {
		return fmt.Errorf("cannot repack zip file %s: %w", output, err)
	}
	return nil
}

// 解压 ZIP 文件到目标目录
func unzip(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fPath := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.Create(fPath)
		if err != nil {
			return err
		}
		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 处理目录中的图片文件
func processImages(dir string, quality int) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// 判断是否为图片文件
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg":
			log.Printf("compress %s", path)
			err := compressJPEG(path, quality)
			if err != nil {
				return fmt.Errorf("cannot compress %s: %w", path, err)
			}
		case ".png":
			log.Printf("compress %s", path)
			err := compress2Jpeg(path, quality)
			if err != nil {
				return fmt.Errorf("cannot compress %s: %w", path, err)
			}
		}
		return nil
	})
}

func compressJPEG(filePath string, quality int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 保存为JPEG（只降低质量，不改变分辨率）
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
}
func compress2Jpeg(filePath string, quality int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return err
	}
	if _, ok := img.(*image.NRGBA); ok { // has alpha channel
		img = fillBackground(img, color.White)
	}
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
}

// 将目录重新打包为 ZIP 文件
func zipDir(srcDir, zipPath string) error {
	of, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer of.Close()

	w := zip.NewWriter(of)
	defer w.Close()

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 处理目录
		if info.IsDir() {
			_, err = w.Create(relPath + "/")
			return err
		}

		// 添加文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := w.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}

// 填充透明背景为指定颜色
func fillBackground(img image.Image, bgColor color.Color) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	// 填充背景
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, bgColor)
		}
	}
	// 将原始图片绘制到目标图像上
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			original := img.At(x, y)
			dst.Set(x, y, original)
		}
	}
	return dst
}
