package download

import (
	"fmt"
	"log"
	"os"
)

func init() {
	err := os.MkdirAll(dirPath, 777)
	if err != nil {
		log.Println("Error creating directory!error: ", err, "path: ", dirPath)
		fmt.Println("创建用于保存下载文件的文件夹时出现错误！请重试！path: ", dirPath)
		return
	}
}
