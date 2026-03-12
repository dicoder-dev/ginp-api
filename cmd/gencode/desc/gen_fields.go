package desc

import (
	"fmt"
	"ginp-api/internal/gen"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GenFields 生成实体常量 - 自动扫描所有实体文件
func GenFields() {
	// 扫描所有实体文件
	entities := ScanExistingEntities()

	if len(entities) == 0 {
		fmt.Println("未找到任何实体文件")
		return
	}

	fmt.Printf("找到 %d 个实体，正在生成字段常量...\n", len(entities))

	// 为每个实体生成字段常量
	for _, entityName := range entities {
		genFieldsForEntity(entityName)
	}

	fmt.Println("字段常量生成完成")
}

// genFieldsForEntity 为单个实体生成字段常量
func genFieldsForEntity(entityName string) {
	// 获取实体文件路径
	entityDir := GetDirEntidy()
	lineName := gen.NameToLine(entityName)
	entityFilePath := filepath.Join(entityDir, lineName+".e.go")

	// 读取实体文件内容
	content, err := os.ReadFile(entityFilePath)
	if err != nil {
		fmt.Printf("读取实体文件失败: %s, 错误: %v\n", entityFilePath, err)
		return
	}

	// 自动检测父级目录
	parentDir := detectEntityParentDirForFields(entityName)

	// 解析字段名
	fields := parseEntityFields(string(content))
	if len(fields) == 0 {
		fmt.Printf("实体 %s 没有找到任何字段\n", entityName)
		return
	}

	// 生成字段常量内容
	packageName := "m" + gen.NameToAllSmall(entityName)
	contentResult := "package " + packageName + " \n\n"

	for _, fieldName := range fields {
		constName := "Field" + fieldName
		fieldNameLower := gen.NameToLine(fieldName)

		// 统一处理时间戳字段
		if fieldName == "CreatedAt" {
			contentResult += fmt.Sprintf("const %s = \"%s\"\n\n", constName, "created_at")
		} else if fieldName == "UpdatedAt" {
			contentResult += fmt.Sprintf("const %s = \"%s\"\n\n", constName, "updated_at")
		} else if fieldName == "DeletedAt" {
			contentResult += fmt.Sprintf("const %s = \"%s\"\n\n", constName, "deleted_at")
		} else {
			contentResult += fmt.Sprintf("const %s = \"%s\"\n\n", constName, fieldNameLower)
		}
	}

	// 写入文件（传入父级目录）
	filePath := PathFields(lineName, parentDir)
	err = os.WriteFile(filePath, []byte(contentResult), 0644)
	if err != nil {
		fmt.Printf("写入文件失败: %s, 错误: %v\n", filePath, err)
		return
	}
	fmt.Printf("实体 %s 的字段常量已生成: %s\n", entityName, filePath)
}

// detectEntityParentDirForFields 检测实体文件所在的父级目录（用于字段常量生成）
func detectEntityParentDirForFields(entityName string) string {
	// 首先尝试从实体配置中获取 FatherFolderName
	fatherFolderName := getEntityFatherFolderName(entityName)
	if fatherFolderName != "" {
		return fatherFolderName
	}

	// 如果没有获取到 FatherFolderName，返回空字符串（表示无父级目录）
	return ""
}

// getEntityFatherFolderName 从实体文件中获取 FatherFolderName 配置
func getEntityFatherFolderName(entityName string) string {
	// 获取实体文件路径
	entityDir := GetDirEntidy()
	lineName := gen.NameToLine(entityName)
	entityFilePath := filepath.Join(entityDir, lineName+".e.go")

	// 读取实体文件内容
	content, err := os.ReadFile(entityFilePath)
	if err != nil {
		return ""
	}

	// 查找 FatherFolderName 配置
	// 匹配模式: FatherFolderName: "xxx"
	fileContent := string(content)
	prefix := "FatherFolderName:"
	startIdx := strings.Index(fileContent, prefix)
	if startIdx == -1 {
		return ""
	}

	// 从匹配位置开始查找引号
	contentAfterPrefix := fileContent[startIdx+len(prefix):]
	// 跳过空格
	contentAfterPrefix = strings.TrimLeft(contentAfterPrefix, " \t")
	// 查找引号
	if !strings.HasPrefix(contentAfterPrefix, "\"") {
		return ""
	}

	// 提取引号内的内容
	endIdx := strings.Index(contentAfterPrefix[1:], "\"")
	if endIdx == -1 {
		return ""
	}

	fatherFolderName := contentAfterPrefix[1 : endIdx+1]
	return fatherFolderName
}

// parseEntityFields 从实体文件内容中解析字段名
func parseEntityFields(content string) []string {
	// 匹配 struct 定义中的字段
	// 匹配 type Xxx struct { 后面的所有字段
	re := regexp.MustCompile(`(?s)type\s+\w+\s+struct\s*\{([^}]+)\}`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return []string{}
	}

	structBody := matches[1]

	// 匹配每个字段: FieldName Type `tag`
	// 使用 (?m) 开启多行模式
	fieldRe := regexp.MustCompile(`(?m)^\s*(\w+)\s+\w+`)
	fieldMatches := fieldRe.FindAllStringSubmatch(structBody, -1)

	fields := []string{}
	seen := make(map[string]bool)
	for _, match := range fieldMatches {
		if len(match) >= 2 {
			fieldName := match[1]
			// 跳过嵌入类型（gorm.Model 等）
			if fieldName != "Model" && fieldName != "Time" {
				if !seen[fieldName] {
					seen[fieldName] = true
					fields = append(fields, fieldName)
				}
			}
		}
	}

	// 额外处理 gorm.Model 嵌入类型
	if strings.Contains(structBody, "gorm.Model") || strings.Contains(structBody, "model") {
		// 添加 ID, CreatedAt, UpdatedAt
		if !seen["ID"] {
			fields = append([]string{"ID"}, fields...)
		}
		if !seen["CreatedAt"] {
			fields = append(fields, "CreatedAt")
		}
		if !seen["UpdatedAt"] {
			fields = append(fields, "UpdatedAt")
		}
	}

	return fields
}
