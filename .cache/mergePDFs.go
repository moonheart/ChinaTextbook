package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	dirPath := "." // 当前目录
	mergeSplitPDFsInDirectoryRecursive(dirPath)
}

func mergeSplitPDFsInDirectoryRecursive(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	// 先处理当前目录下的PDF分卷文件
	mergeSplitPDFsInDirectory(dirPath)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		dirName := file.Name()
		if strings.HasSuffix(dirName, ".pdf_merge_folder") {
			// 合并该目录下的PDF文件
			mergeSplitPDFsInDirectory(filepath.Join(dirPath, dirName))

			// 获取合并后的PDF文件名(去掉.pdf_merge_folder后缀)
			baseName := strings.TrimSuffix(dirName, ".pdf_merge_folder") + ".pdf"
			targetPath := filepath.Join(dirPath, baseName)

			// 重命名合并后的文件到目标位置
			if err := os.Rename(filepath.Join(dirPath, dirName, baseName), targetPath); err != nil {
				panic(err)
			}

			// 删除.pdf_merge_folder目录
			if err := os.RemoveAll(filepath.Join(dirPath, dirName)); err != nil {
				panic(err)
			}
		} else {
			// 递归处理子目录
			mergeSplitPDFsInDirectoryRecursive(filepath.Join(dirPath, dirName))
		}
	}
}

func mergeSplitPDFsInDirectory(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	splitFiles := make(map[string][]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if strings.Contains(fileName, ".pdf.") {
			baseName := strings.Split(fileName, ".pdf.")[0] + ".pdf"
			splitFiles[baseName] = append(splitFiles[baseName], filepath.Join(dirPath, fileName))
		}
	}

	for baseName, parts := range splitFiles {
		sort.Strings(parts) // 确保文件顺序正确
		mergeFiles(filepath.Join(dirPath, baseName), parts)
	}
}

func mergeFiles(baseName string, parts []string) {
	mergedFile, err := os.Create(baseName)
	if err != nil {
		panic(err)
	}
	defer mergedFile.Close()

	for _, part := range parts {
		data, err := os.ReadFile(part)
		if err != nil {
			panic(err)
		}
		_, err = mergedFile.Write(data)
		if err != nil {
			panic(err)
		}
		os.Remove(part) // 合并后删除分割文件
	}
}
