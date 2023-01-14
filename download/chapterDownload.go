package download

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"sese/getMangaInfo"
	"strconv"
	"sync"
	"sync/atomic"
)

const dirPath = "./comic_download"

func DownloadAChapter(pic []getMangaInfo.ComicPic, info getMangaInfo.ComicInfo, chapter string, class string) {
	name := info.Name
	path := dirPath + "/" + name + "_" + info.PathWord + "/" + class + "/" + chapter + "/"
	err := os.MkdirAll(path, 777)
	if err != nil {
		log.Println("Error creating directory!error: ", err, "path: ", path)
		fmt.Println("创建用于保存下载文件的文件夹时出现错误！请重试！path: ", path)
		return
	}
	var wg sync.WaitGroup
	var failed atomic.Int32
	for i, v := range pic {
		wg.Add(1)
		go func(index int, comicPic getMangaInfo.ComicPic) {
			defer wg.Done()
			var get *http.Response
			var err error
			var success bool
			for i := 0; i < 3; i++ {
				get, err = http.Get(comicPic.Url)
				if err != nil {
					log.Println("获取网络图片时出错，name:", name, "chapter:", chapter, "page:", index, "error:", err)
					failed.Add(1)
					return
				}
				if get.StatusCode == http.StatusOK {
					success = true
					break
				}
			}
			defer get.Body.Close()
			if !success {
				log.Println("获取与资源之间的连接失败，name:", name, "chapter:", chapter, "page:", index, "error:", err)
				failed.Add(1)
				return
			}
			bar := progressbar.DefaultBytes(
				get.ContentLength,
				"正在下载"+chapter+"第"+strconv.Itoa(index)+"页",
			)
			f, err := os.OpenFile(path+strconv.Itoa(index)+".png", os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("创建并打开本地文件时出错，name:", name, "chapter:", chapter, "page:", index, "error:", err)
				failed.Add(1)
				return
			}
			defer f.Close()
			_, err = io.Copy(io.MultiWriter(f, bar), get.Body)
			if err != nil {
				log.Println("读取网络图片字节流时出错，name:", name, "chapter:", chapter, "page:", index, "error:", err)
				failed.Add(1)
				return
			}
		}(i+1, v)
	}
	wg.Wait()
	if failed.Load() != 0 {
		fmt.Println(name+chapter+"下载结束，共", len(pic), "页，其中有", failed.Load(), "页下载失败")
	} else {
		fmt.Println(name+chapter+"下载成功!共", len(pic), "页")
	}
}
