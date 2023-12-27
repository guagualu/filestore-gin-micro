package main

import (
	"fileStore/internel/pkg/util"
	"fmt"
)

func main() {
	srcPath := "./static/mp/1234567"
	destPath := "../../tmp/" + "f39aa2949db19311413d567c656103367f60c84c"
	////cd ./static/mp/1234567 && ls | sort -n | xargs cat > ../../tmp/f39aa2949db19311413d567c656103367f60c84c
	//cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	//cmd = fmt.Sprintf("merge.sh")
	//cd./static/mp/1234567; Get-ChildItem -Name | Sort-Object -Numeric | ForEach-Object { Get-Content "./static/mp/1234567/$_" -Encoding Byte } >../../tmp/f39aa2949db19311413d567c656103367f60c84c
	out, Err, err := util.ExecWinShell("merge.sh", fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath))
	if err != nil {
		fmt.Println("分块文件合并失败:", err)
		fmt.Println("分块文件合并失败:", Err)
		return
	}
	fmt.Println(out, Err)

}
