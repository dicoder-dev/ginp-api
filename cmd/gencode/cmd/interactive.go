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
		fmt.Print("请输入选项编号: ")

		choice := readInput()

		switch choice {
		case "1":
			handleGenCrud()
		case "2":
			handleAddApi()
		case "3":
			handleGenFields()
		case "4":
			handleRemoveCrud()
		case "5":
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
		fmt.Println("请选择目标目录（输入编号）：")
		for i, dir := range controllerDirs {
			fmt.Printf("%d. %s\n", i+1, dir)
		}
		fmt.Println(fmt.Sprintf("%d. [创建新目录]", len(controllerDirs)+1))
		fmt.Print("请输入选项编号: ")

		dirChoice := readInput()
		dirChoice = strings.TrimSpace(dirChoice)

		dirIndex, err := strconv.Atoi(dirChoice)
		if err != nil || dirIndex < 1 || dirIndex > len(controllerDirs)+1 {
			fmt.Println("无效的选项")
			return
		}

		var parentDir string
		if dirIndex <= len(controllerDirs) {
			parentDir = controllerDirs[dirIndex-1]
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
			fmt.Println("请选择子目录（输入编号）：")
			for i, dir := range subDirs {
				fmt.Printf("%d. %s\n", i+1, dir)
			}
			fmt.Println(fmt.Sprintf("%d. [在父级目录下新建]", len(subDirs)+1))
			fmt.Print("请输入选项编号: ")

			subDirChoice := readInput()
			subDirChoice = strings.TrimSpace(subDirChoice)

			subDirIndex, err := strconv.Atoi(subDirChoice)
			if err != nil || subDirIndex < 1 || subDirIndex > len(subDirs)+1 {
				fmt.Println("无效的选项")
				return
			}

			if subDirIndex <= len(subDirs) {
				// 选择已有子目录
				finalDirPath = parentDir + "/c" + subDirs[subDirIndex-1]
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
			fmt.Println("请选择存放位置（输入编号）：")
			fmt.Println("1. 在父级目录下新建")
			if dirIndex <= len(controllerDirs) {
				fmt.Println("2. 存放到已有的 c" + parentDir + " 目录（如果存在）")
			}
			fmt.Print("请输入选项编号: ")

			locationChoice := readInput()
			locationChoice = strings.TrimSpace(locationChoice)

			if locationChoice == "2" && dirIndex <= len(controllerDirs) {
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
	}
}

// handleGenCrud 处理生成实体 CRUD
func handleGenCrud() {
	fmt.Println("")
	fmt.Println("=== 生成实体 CRUD 代码 ===")

	// 扫描已存在的实体
	existingEntities := desc.ScanExistingEntities()

	// 显示选项
	fmt.Println("请选择操作：")
	fmt.Println("1. 创建新实体并生成 CRUD")
	if len(existingEntities) > 0 {
		fmt.Println("2. 为已存在实体生成 CRUD")
	}
	fmt.Print("请输入选项编号: ")

	choice := readInput()
	choice = strings.TrimSpace(choice)

	if choice == "1" || (choice == "2" && len(existingEntities) == 0) {
		// 创建新实体
		handleCreateNewEntity()
	} else if choice == "2" && len(existingEntities) > 0 {
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
		fmt.Println("请选择父级目录（直接回车使用默认目录）：")
		for i, dir := range controllerDirs {
			fmt.Printf("%d. %s\n", i+1, dir)
		}
		fmt.Printf("0. 不使用父级目录\n")
		fmt.Print("请输入选项编号: ")

		dirChoice := readInput()
		dirChoice = strings.TrimSpace(dirChoice)

		if dirChoice != "" {
			dirIndex, err := strconv.Atoi(dirChoice)
			if err == nil && dirIndex >= 0 && dirIndex <= len(controllerDirs) {
				if dirIndex > 0 {
					parentDir = controllerDirs[dirIndex-1]
				}
			}
		}
	}

	// 生成 CRUD
	entities := []string{entityName}
	desc.GenBatchCrudWithParent(entities, parentDir)
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
		fmt.Println("请选择父级目录（直接回车使用默认目录）：")
		for i, dir := range controllerDirs {
			fmt.Printf("%d. %s\n", i+1, dir)
		}
		fmt.Printf("0. 不使用父级目录\n")
		fmt.Print("请输入选项编号: ")

		dirChoice := readInput()
		dirChoice = strings.TrimSpace(dirChoice)

		if dirChoice != "" {
			dirIndex, err := strconv.Atoi(dirChoice)
			if err == nil && dirIndex >= 0 && dirIndex <= len(controllerDirs) {
				if dirIndex > 0 {
					parentDir = controllerDirs[dirIndex-1]
				}
			}
		}
	}

	// 生成 CRUD
	desc.GenBatchCrudWithParent(selectedEntities, parentDir)
}

// handleGenFields 处理生成实体字段常量
func handleGenFields() {
	fmt.Println("")
	fmt.Println("=== 生成实体字段常量 ===")
	fmt.Println("正在生成字段常量...")

	desc.GenFields()

	fmt.Println("字段常量生成完成")
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
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("操作已取消")
		return
	}

	// 删除每个实体
	for _, entityName := range selectedEntities {
		parentDir := entityParentMap[entityName]
		desc.RemoveBatchCrudWithParent([]string{entityName}, parentDir)
	}
}

// detectEntityParentDir 检测实体文件所在的父级目录
func detectEntityParentDir(entityName string) string {
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
