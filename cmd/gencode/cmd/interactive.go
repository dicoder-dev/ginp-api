package cmd

import (
	"bufio"
	"fmt"
	"ginp-api/cmd/gencode/desc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 创建一个共享的 reader
var sharedReader *bufio.Reader

// initSharedReader 初始化共享的 reader
func initSharedReader() {
	if sharedReader == nil {
		sharedReader = bufio.NewReader(os.Stdin)
	}
}

// readInput 读取用户输入
func readInput() string {
	initSharedReader()

	// 刷新标准输出缓冲区，确保提示信息显示
	os.Stdout.Sync()

	input, err := sharedReader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}

// readChoiceWithDefault 读取用户选择，回车默认选择第一个选项
// 返回选择的索引（从0开始）
func readChoiceWithDefault(maxChoice int) int {
	input := readInput()
	if input == "" {
		// 回车默认选择第一个
		return 0
	}
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > maxChoice {
		// 无效输入返回 -1
		return -1
	}
	return choice - 1
}

// waitForEnter 提示用户按回车继续
func waitForEnter() {
	fmt.Print("\n按回车键继续...")
	readInput()
}

// ShowMainMenu 显示主菜单并处理用户选择
func ShowMainMenu() {
	for {
		fmt.Println("")
		fmt.Println("=== GAPI 代码生成工具 ===")
		fmt.Println("请选择操作：")
		fmt.Println("1. 生成实体 CRUD 代码")
		fmt.Println("2. 新增 API 接口控制器")
		fmt.Println("3. 生成实体字段常量")
		fmt.Println("4. 删除实体 CRUD 代码")
		fmt.Println("5. 退出")
		fmt.Print("请输入选项编号（直接回车默认选择 1）: ")

		choice := readChoiceWithDefault(5)

		switch choice {
		case 0:
			handleGenCrud()
		case 1:
			handleAddApi()
		case 2:
			handleGenFields()
		case 3:
			handleRemoveCrud()
		case 4:
			fmt.Println("退出程序")
			return
		default:
			fmt.Println("无效的选项，请重新选择")
		}
	}
}

// handleAddApi 处理新增 API 接口
func handleAddApi() {
	// 扫描控制器目录
	controllerDirs := desc.ScanControllerDirs()

	fmt.Println("")
	fmt.Println("=== 新增 API 接口控制器 ===")

	// 显示现有目录列表
	if len(controllerDirs) > 0 {
		fmt.Println("请选择目标目录（输入编号，直接回车默认选择 1）：")
		for i, dir := range controllerDirs {
			fmt.Printf("%d. %s\n", i+1, dir)
		}
		fmt.Println(fmt.Sprintf("%d. [创建新目录]", len(controllerDirs)+1))
		fmt.Print("请输入选项编号: ")

		dirChoice := readChoiceWithDefault(len(controllerDirs) + 1)

		if dirChoice == -1 {
			fmt.Println("无效的选项")
			return
		}

		var parentDir string
		if dirChoice < len(controllerDirs) {
			parentDir = controllerDirs[dirChoice]
		} else {
			// 创建新目录
			fmt.Print("请输入新目录名称（如 user）: ")
			parentDir = readInput()
			parentDir = strings.TrimSpace(parentDir)
			if parentDir == "" {
				fmt.Println("目录名称不能为空")
				return
			}
		}

		// 检查父级目录下是否有子目录（以c开头的）
		subDirs := desc.ScanSubDirs(parentDir)
		var finalDirPath string

		if len(subDirs) > 0 {
			// 有子目录，让用户选择
			fmt.Println("请选择子目录（输入编号，直接回车默认选择 1）：")
			for i, dir := range subDirs {
				fmt.Printf("%d. %s\n", i+1, dir)
			}
			fmt.Println(fmt.Sprintf("%d. [在父级目录下新建]", len(subDirs)+1))
			fmt.Print("请输入选项编号: ")

			subDirChoice := readChoiceWithDefault(len(subDirs) + 1)

			if subDirChoice == -1 {
				fmt.Println("无效的选项")
				return
			}

			if subDirChoice < len(subDirs) {
				// 选择已有子目录
				finalDirPath = parentDir + "/c" + subDirs[subDirChoice]
			} else {
				// 在父级目录下新建（新建的文件夹必须以c开头）
				fmt.Print("请输入新子目录名称（如 user）: ")
				newSubDir := readInput()
				newSubDir = strings.TrimSpace(newSubDir)
				if newSubDir == "" {
					fmt.Println("目录名称不能为空")
					return
				}
				finalDirPath = parentDir + "/c" + newSubDir
			}
		} else {
			// 没有子目录，直接在父级目录下新建（新建的文件夹必须以c开头）
			// 如果用户选择了已有目录，也需要添加到c开头的目录下
			maxLocation := 1
			if dirChoice < len(controllerDirs) {
				maxLocation = 2
			}
			fmt.Println("请选择存放位置（输入编号，直接回车默认选择 1）：")
			fmt.Println("1. 在父级目录下新建")
			if dirChoice < len(controllerDirs) {
				fmt.Println("2. 存放到已有的 c" + parentDir + " 目录（如果存在）")
			}
			fmt.Print("请输入选项编号: ")

			locationChoice := readChoiceWithDefault(maxLocation)

			if locationChoice == -1 {
				fmt.Println("无效的选项")
				return
			}

			if locationChoice == 1 && dirChoice < len(controllerDirs) {
				// 存放到已有的 c + parentDir 目录
				// 先检查该目录是否存在
				existingDir := "c" + parentDir
				controllerBase := desc.GetDirController()
				existingPath := filepath.Join(controllerBase, existingDir)
				if _, err := os.Stat(existingPath); err == nil {
					finalDirPath = parentDir + "/" + existingDir
				} else {
					// 目录不存在，新建
					fmt.Printf("目录 %s 不存在，将创建该目录\n", existingDir)
					finalDirPath = parentDir + "/c" + parentDir
				}
			} else {
				// 在父级目录下新建（新建的文件夹必须以c开头）
				fmt.Print("请输入新目录名称（如 user）: ")
				newDir := readInput()
				newDir = strings.TrimSpace(newDir)
				if newDir == "" {
					fmt.Println("目录名称不能为空")
					return
				}
				finalDirPath = parentDir + "/c" + newDir
			}
		}

		// 获取 API 名称
		fmt.Print("请输入 API 名称（大驼峰命名，如 GetUserInfo）: ")
		apiName := readInput()
		apiName = strings.TrimSpace(apiName)

		if apiName == "" {
			fmt.Println("API 名称不能为空")
			return
		}

		// 调用现有函数生成 API
		desc.GenAddApiWithParams(apiName, finalDirPath)
		waitForEnter()
	} else {
		// 没有现有目录，引导创建
		fmt.Println("当前没有控制器目录")
		fmt.Print("请输入新目录名称（如 user）: ")

		dirPath := readInput()
		dirPath = strings.TrimSpace(dirPath)

		if dirPath == "" {
			fmt.Println("目录名称不能为空")
			return
		}

		// 新建的接口文件夹必须以c开头
		dirPath = "c" + dirPath

		// 获取 API 名称
		fmt.Print("请输入 API 名称（大驼峰命名，如 GetUserInfo）: ")
		apiName := readInput()
		apiName = strings.TrimSpace(apiName)

		if apiName == "" {
			fmt.Println("API 名称不能为空")
			return
		}

		// 调用现有函数生成 API
		desc.GenAddApiWithParams(apiName, dirPath)
		waitForEnter()
	}
}

// handleGenCrud 处理生成实体 CRUD
func handleGenCrud() {
	fmt.Println("")
	fmt.Println("=== 生成实体 CRUD 代码 ===")

	// 扫描已存在的实体
	existingEntities := desc.ScanExistingEntities()

	// 显示选项
	fmt.Println("请选择操作（直接回车默认选择 1）：")
	fmt.Println("1. 创建新实体并生成 CRUD")
	if len(existingEntities) > 0 {
		fmt.Println("2. 为已存在实体生成 CRUD")
	}
	fmt.Print("请输入选项编号: ")

	maxChoice := 1
	if len(existingEntities) > 0 {
		maxChoice = 2
	}
	choice := readChoiceWithDefault(maxChoice)

	if choice == 0 || (choice == 1 && len(existingEntities) == 0) {
		// 创建新实体
		handleCreateNewEntity()
	} else if choice == 1 && len(existingEntities) > 0 {
		// 为已存在实体生成 CRUD
		handleGenCrudForExisting(existingEntities)
	} else {
		fmt.Println("无效的选项")
	}
}

// handleCreateNewEntity 处理创建新实体
func handleCreateNewEntity() {
	fmt.Println("")
	fmt.Println("=== 创建新实体并生成 CRUD ===")

	// 获取实体名称
	fmt.Print("请输入实体名称（大驼峰命名，如 UserGroup）: ")
	entityName := readInput()
	entityName = strings.TrimSpace(entityName)

	if entityName == "" {
		fmt.Println("实体名称不能为空")
		return
	}

	// 获取父级目录
	controllerDirs := desc.ScanControllerDirs()
	var parentDir string

	if len(controllerDirs) > 0 {
		fmt.Println("请选择父级目录（直接回车默认选择 0 不使用父级目录）：")
		for i, dir := range controllerDirs {
			fmt.Printf("%d. %s\n", i+1, dir)
		}
		fmt.Printf("0. 不使用父级目录\n")
		fmt.Print("请输入选项编号: ")

		// 选项范围是 1 到 len(controllerDirs)，0 表示不使用父级目录
		maxChoice := len(controllerDirs)
		dirChoice := readChoiceWithDefault(maxChoice)

		if dirChoice == -1 {
			fmt.Println("无效的选项")
			return
		}

		fmt.Printf("DEBUG: len(controllerDirs)=%d, dirChoice=%d\n", len(controllerDirs), dirChoice)

		// dirChoice 是从 0 开始的索引（0 表示不使用父级目录，1 表示第一个目录，以此类推）
		if dirChoice > 0 && dirChoice <= len(controllerDirs) {
			parentDir = controllerDirs[dirChoice]
		}
	}

	// 生成 CRUD
	entities := []string{entityName}
	desc.GenBatchCrudWithParent(entities, parentDir)
	waitForEnter()
}

// handleGenCrudForExisting 处理为已存在实体生成 CRUD
func handleGenCrudForExisting(existingEntities []string) {
	fmt.Println("")
	fmt.Println("=== 为已存在实体生成 CRUD ===")

	// 显示已存在的实体
	fmt.Println("已存在的实体：")
	for i, entity := range existingEntities {
		fmt.Printf("%d. %s\n", i+1, entity)
	}
	fmt.Print("请输入实体编号（多个用逗号分隔，如 1,2,3）: ")

	choice := readInput()
	choice = strings.TrimSpace(choice)

	if choice == "" {
		fmt.Println("输入不能为空")
		return
	}

	// 解析选择
	selectedEntities := []string{}
	choices := strings.Split(choice, ",")

	for _, c := range choices {
		c = strings.TrimSpace(c)
		index, err := strconv.Atoi(c)
		if err == nil && index >= 1 && index <= len(existingEntities) {
			selectedEntities = append(selectedEntities, existingEntities[index-1])
		}
	}

	if len(selectedEntities) == 0 {
		fmt.Println("未选择任何实体")
		return
	}

	// 获取父级目录
	controllerDirs := desc.ScanControllerDirs()
	var parentDir string

	if len(controllerDirs) > 0 {
		fmt.Println("请选择父级目录（直接回车默认选择 0 不使用父级目录）：")
		for i, dir := range controllerDirs {
			fmt.Printf("%d. %s\n", i+1, dir)
		}
		fmt.Printf("0. 不使用父级目录\n")
		fmt.Print("请输入选项编号: ")

		// 选项范围是 1 到 len(controllerDirs)，0 表示不使用父级目录
		maxChoice := len(controllerDirs)
		dirChoice := readChoiceWithDefault(maxChoice)

		if dirChoice == -1 {
			fmt.Println("无效的选项")
			return
		}

		fmt.Printf("DEBUG: len(controllerDirs)=%d, dirChoice=%d\n", len(controllerDirs), dirChoice)

		// dirChoice 是从 0 开始的索引（0 表示不使用父级目录，1 表示第一个目录，以此类推）
		if dirChoice > 0 && dirChoice <= len(controllerDirs) {
			parentDir = controllerDirs[dirChoice]
		}
	}

	// 生成 CRUD
	desc.GenBatchCrudWithParent(selectedEntities, parentDir)
	waitForEnter()
}

// handleGenFields 处理生成实体字段常量
func handleGenFields() {
	fmt.Println("")
	fmt.Println("=== 生成实体字段常量 ===")
	fmt.Println("正在生成字段常量...")

	desc.GenFields()

	fmt.Println("字段常量生成完成")
	waitForEnter()
}

// handleRemoveCrud 处理删除实体 CRUD
func handleRemoveCrud() {
	fmt.Println("")
	fmt.Println("=== 删除实体 CRUD 代码 ===")

	// 扫描已存在的实体
	existingEntities := desc.ScanExistingEntities()

	if len(existingEntities) == 0 {
		fmt.Println("当前没有已存在的实体")
		return
	}

	// 显示已存在的实体
	fmt.Println("已存在的实体：")
	for i, entity := range existingEntities {
		fmt.Printf("%d. %s\n", i+1, entity)
	}
	fmt.Print("请输入实体编号（多个用逗号分隔，如 1,2,3）: ")

	choice := readInput()
	choice = strings.TrimSpace(choice)

	if choice == "" {
		fmt.Println("输入不能为空")
		return
	}

	// 解析选择
	selectedEntities := []string{}
	choices := strings.Split(choice, ",")

	for _, c := range choices {
		c = strings.TrimSpace(c)
		index, err := strconv.Atoi(c)
		if err == nil && index >= 1 && index <= len(existingEntities) {
			selectedEntities = append(selectedEntities, existingEntities[index-1])
		}
	}

	if len(selectedEntities) == 0 {
		fmt.Println("未选择任何实体")
		return
	}

	// 自动检测每个实体的正确父级目录
	entityParentMap := make(map[string]string) // entityName -> parentDir
	for _, entityName := range selectedEntities {
		parentDir := detectEntityParentDir(entityName)
		entityParentMap[entityName] = parentDir
		if parentDir != "" {
			fmt.Printf("实体 %s 位于父级目录: %s\n", entityName, parentDir)
		} else {
			fmt.Printf("实体 %s 无父级目录\n", entityName)
		}
	}

	// 确认删除操作
	fmt.Printf("\n确认删除以下实体的 CRUD 代码吗？(%s) [y/N]: ", strings.Join(selectedEntities, ", "))
	confirm := readInput()
	// 按回车或输入 y/yes 确认，其他输入取消
	if confirm != "" && strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("操作已取消")
		return
	}

	// 删除每个实体
	for _, entityName := range selectedEntities {
		parentDir := entityParentMap[entityName]
		desc.RemoveBatchCrudWithParent([]string{entityName}, parentDir)
	}
	waitForEnter()
}

// detectEntityParentDir 检测实体文件所在的父级目录
func detectEntityParentDir(entityName string) string {
	// 首先尝试从实体配置中获取 FatherFolderName
	fatherFolderName := getEntityFatherFolderName(entityName)
	if fatherFolderName != "" {
		// 使用 FatherFolderName + c + 实体名（小写首字母）构建路径
		lineName := strings.ToLower(entityName)
		if len(lineName) > 0 {
			lineName = strings.ToLower(string(lineName[0])) + lineName[1:]
		}
		lineName = strings.ReplaceAll(lineName, "_", "")

		// 构建完整路径: FatherFolderName/cUserTest
		controllerBase := desc.GetDirController()
		entityDir := filepath.Join(controllerBase, fatherFolderName, "c"+lineName)

		if _, err := os.Stat(entityDir); err == nil {
			return fatherFolderName
		}
	}

	// 如果没有获取到 FatherFolderName，使用原来的检测逻辑
	lineName := strings.ToLower(entityName)
	if len(lineName) > 0 {
		lineName = strings.ToLower(string(lineName[0])) + lineName[1:]
	}
	lineName = strings.ReplaceAll(lineName, "_", "")

	// 检查 controller 目录下的所有子目录
	controllerBase := desc.GetDirController()
	entries, _ := os.ReadDir(controllerBase)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// 检查这个目录下是否有 c + lineName 的子目录
		entityDir := filepath.Join(controllerBase, entry.Name(), "c"+lineName)
		if _, err := os.Stat(entityDir); err == nil {
			return entry.Name()
		}
	}

	// 检查顶层是否有 c + lineName 目录（无父级目录的情况）
	topLevelDir := filepath.Join(controllerBase, "c"+lineName)
	if _, err := os.Stat(topLevelDir); err == nil {
		return ""
	}

	return ""
}

// getEntityFatherFolderName 从实体文件中获取 FatherFolderName 配置
func getEntityFatherFolderName(entityName string) string {
	// 获取实体文件路径
	entityDir := desc.GetDirEntidy()
	lineName := strings.ToLower(entityName)
	if len(lineName) > 0 {
		lineName = strings.ToLower(string(lineName[0])) + lineName[1:]
	}
	lineName = strings.ReplaceAll(lineName, "_", "")
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
