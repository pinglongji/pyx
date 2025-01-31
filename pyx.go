package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// 缓存依赖项的路径
var depsCache string

func init() {
	// 将缓存路径初始化为几个可能的位置
	if home := os.Getenv("HOME"); home != "" {
		depsCache = filepath.Join(home, ".pyx-cache")
		return
	}
	if user, err := user.Current(); user != nil && err == nil && user.HomeDir != "" {
		depsCache = filepath.Join(user.HomeDir, ".pyx-cache")
		return
	}
	depsCache = filepath.Join(os.TempDir(), "pyx-cache")
}

func main() {
	// 确保docker可用
	if err := checkDocker(); err != nil {
		log.Fatalf("Failed to check docker installation: %v.", err)
	}
	// 检查所有必需的镜像是否可用
	image := "x86_64_pyx"
	found, err := checkDockerImage(image)
	switch {
	case err != nil:
		log.Fatalf("Failed to check docker image availability: %v.", err)
	case !found:
		log.Fatalf("image not found: %v", image)
	default:
		fmt.Println("image found.")
	}

	folder, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to retrieve the working directory: %v.", err)
	}

	err = compile(image, folder)
}

// 检查是否可以找到一个安装的docker并运行正常。
func checkDocker() error {
	fmt.Println("Checking docker installation...")
	if err := run(exec.Command("docker", "version")); err != nil {
		return err
	}
	return nil
}

// 同步执行命令，将其输出重定向到标准输出。
func run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// 检查所需的docker镜像是否在本地可用。
func checkDockerImage(image string) (bool, error) {
	fmt.Printf("Checking for required docker image %s... \n", image)
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return false, err
	}
	return bytes.Contains(out, []byte(image)), nil
}

// 根据给定的构建规范，使用特定的docker镜像编译请求的包。
func compile(image string, folder string) error {
	locals, mounts := []string{}, []string{}
	// 内部依赖包的位置
	path := filepath.Join(folder, ".deps.txt")
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("The deps.txt not exist in the %v, dependency is not loaded", folder)
	} else {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			target := scanner.Text()
			target = strings.TrimSpace(target)
			// 跳过空行
			if target == "" {
				continue
			}
			locals = append(locals, target)
			mounts = append(mounts, filepath.Join("/private-cache", filepath.Base(target)))
		}

	}

	// 组装并运行交叉编译命令
	fmt.Printf("starting compiling ...\n")
	args := []string{
		"run", "--rm",
		"-v", folder + ":/build",
		"-v", depsCache + ":/deps-cache:rw",
		"-e", "PYTHONPATH=/private-cache",
		"-e", "LC_ALL=en_US.utf8",
		"-e", "LANG=en_US.utf8",
	}

	for i := 0; i < len(locals); i++ {
		args = append(args, []string{"-v", fmt.Sprintf("%s:%s:ro", locals[i], mounts[i])}...)
	}

	args = append(args, image)

	return run(exec.Command("docker", args...))
}
