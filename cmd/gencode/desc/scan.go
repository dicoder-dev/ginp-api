package desc

import (
	"fmt"
	"ginp-api/internal/gen"
	"os"
	"path/filepath"
	"strings"
)

// ScanExistingEntities 扫描已存在的实体
// 返回实体名称列表（大驼峰格式）
func ScanExistingEntities() []string {
	entityDir := GetDirEntidy()
	entities := []string{}

	// 读取 entity 目录下的所有 .e.go 文件
	files, err := os.ReadDir(entityDir)
	if err != nil {
		fmt.Printf("读取实体目录失败: %s\n", err)
		return entities
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".e.go") {
			// 从文件名提取实体名称，如 user_group.e.go -> UserGroup
			baseName := strings.TrimSuffix(file.Name(), ".e.go")
			entityName := gen.NameToCameBig(baseName)
			entities = append(entities, entityName)
		}
	}

	return entities
}

// ScanControllerDirs 扫描控制器目录
// 返回目录名称列表（不带 c 前缀）
func ScanControllerDirs() []string {
	controllerDir := GetDirController()
	dirs := []string{}

	// 读取 controller 目录下的所有子目录
	entries, err := os.ReadDir(controllerDir)
	if err != nil {
		fmt.Printf("读取控制器目录失败: %s\n", err)
		return dirs
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			// 去掉 c 前缀（如果有）
			if strings.HasPrefix(dirName, "c") {
				dirName = strings.TrimPrefix(dirName, "c")
			}
			dirs = append(dirs, dirName)
		}
	}

	return dirs
}

// GetAllControllerDirsWithPath 获取所有控制器目录的完整路径和名称
// 返回 map[dirName]fullPath
func GetAllControllerDirsWithPath() map[string]string {
	controllerDir := GetDirController()
	dirs := make(map[string]string)

	entries, err := os.ReadDir(controllerDir)
	if err != nil {
		fmt.Printf("读取控制器目录失败: %s\n", err)
		return dirs
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			fullPath := filepath.Join(controllerDir, dirName)
			dirs[dirName] = fullPath
		}
	}

	return dirs
}
