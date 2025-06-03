package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"

	"github.com/disintegration/imaging"
)

func main() {
	// 用 pngpaste 读取剪贴板图片到临时文件
	tmpFile := "clipboard.png"
	cmd := exec.Command("pngpaste", tmpFile)
	if err := cmd.Run(); err != nil {
		fmt.Println("剪贴板没有图片或读取失败:", err)
		return
	}
	defer os.Remove(tmpFile)

	// 打开图片
	file, err := os.Open(tmpFile)
	if err != nil {
		fmt.Println("打开临时图片失败:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("图片解码失败:", err)
		return
	}

	// 这里没有直接获取 DPI 的方法，假设为 96
	wmH := 70

	// 如果图片太窄
	if img.Bounds().Dx() <= 400 {
		return
	}

	// 裁剪图片
	rect := image.Rect(0, wmH, img.Bounds().Dx(), img.Bounds().Dy())
	cropped := imaging.Crop(img, rect)

	// 保存到 buffer
	var buf bytes.Buffer
	err = png.Encode(&buf, cropped)
	if err != nil {
		fmt.Println("图片编码失败:", err)
		return
	}

	// 保存裁剪后的图片到临时文件
	tmpOut := "cropped.png"
	err = os.WriteFile(tmpOut, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("保存裁剪图片失败:", err)
		return
	}
	defer os.Remove(tmpOut)

	// 用 osascript 写回剪贴板
	cmd = exec.Command("osascript", "-e", fmt.Sprintf(`set the clipboard to (read (POSIX file "%s") as «class PNGf»)`, tmpOut))
	if err := cmd.Run(); err != nil {
		fmt.Println("写入剪贴板失败:", err)
		return
	}
	fmt.Println("已裁剪并写回剪贴板")
}
