package main

import (
	"fmt"
	"os"
)

func main() {
	dir := "wind-breaker"
	//files, _ := os.ReadDir(dir)
	//for _, file := range files {
	//	if !file.IsDir() {
	//		chapter := strings.Split(file.Name(), "-")[1]
	//		oldPath := fmt.Sprintf("%s/%s", dir, file.Name())
	//		newPath := fmt.Sprintf("%s/Chapter-%s/%s", dir, chapter, file.Name())
	//		err := os.Rename(oldPath, newPath)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}
	//}
	for i := range 500 {
		err := os.Mkdir(fmt.Sprintf("%s/Chapter-%d", dir, i), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
