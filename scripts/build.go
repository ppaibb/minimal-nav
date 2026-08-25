package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// copyDir 递归复制目录
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

func main() {
	fmt.Println("🚀 开始构建 Minimal Nav 单二进制可执行程序...")

	// 1. 同步 frontend/dist 到 backend/dist
	frontendDist := filepath.Join("frontend", "dist")
	backendDist := filepath.Join("backend", "dist")

	if _, err := os.Stat(frontendDist); os.IsNotExist(err) {
		fmt.Println("❌ 未检测到 frontend/dist，请先执行 npm run build --prefix frontend")
		os.Exit(1)
	}

	_ = os.RemoveAll(backendDist)
	if err := copyDir(frontendDist, backendDist); err != nil {
		fmt.Printf("❌ 复制静态资源失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 前端静态资源已成功同步至 backend/dist")

	// 2. 执行 Go 编译
	outName := "minimal-nav"
	if os.Getenv("GOOS") == "windows" || (os.Getenv("GOOS") == "" && filepath.Separator == '\\') {
		outName += ".exe"
	}

	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", filepath.Join("..", outName), ".")
	cmd.Dir = "backend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Go 编译失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🎉 单二进制构建成功: %s\n", outName)
	fmt.Println("💡 你可以直接双击或运行该文件，即可在本地/服务器直接启动完整导航站！")
}
