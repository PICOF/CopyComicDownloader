package main

import (
	"flag"
	"fmt"
	"sese/download"
	"sese/getMangaInfo"
	"strings"
	"sync"
)

var name *string
var downloadAll *bool
var limit *int
var concurrent *bool

func init() {
	name = flag.String("name", "电锯人", "漫画名字，默认为电锯人")
	limit = flag.Int("limit", 10, "搜索漫画的返回结果条数,默认为十条")
	downloadAll = flag.Bool("all", false, "是否整本下载，开启后无法指定章节")
	concurrent = flag.Bool("concurrent", false, "是否在整本下载时针对每个章节都采用并发下载（速度快，但容易因为请求过多导致问题）")
	flag.Parse()
}

func main() {
	var pic []getMangaInfo.ComicPic
	var typeMap = make(map[int]string)
	res := getMangaInfo.SearchComic(*limit, *name)
	if len(res.Result.List) == 0 {
		fmt.Println("未找到相关漫画！")
		return
	}
	fmt.Print(res.GetSearchInfo())
	fmt.Println("请输入您想下载的漫画的序号")
	var index int
	_, err := fmt.Scan(&index)
	if err != nil {
		fmt.Println("扫描出错，请稍后重试！")
		return
	}
	info := getMangaInfo.GetComicDetails(res.Result.List[index].PathWord)
	for _, v := range info.Build.Type {
		typeMap[v.Id] = v.Name
	}
	if *downloadAll {
		var wg sync.WaitGroup
		for _, v := range info.Groups.Default.Chapters {
			wg.Add(1)
			if *concurrent {
				go func(c getMangaInfo.Chapter) {
					defer wg.Done()
					pic = info.GetChaptersInfo(c.Id)
					download.DownloadAChapter(pic, res.Result.List[index], c.Name, typeMap[c.Type])
				}(v)
			} else {
				pic = info.GetChaptersInfo(v.Id)
				download.DownloadAChapter(pic, res.Result.List[index], v.Name, typeMap[v.Type])
				wg.Done()
			}
		}
		wg.Wait()
		fmt.Println("整本下载已结束！请检查相应文件是否成功下载")
	} else {
		list := info.GetList()
		for k, v := range list {
			fmt.Println(k)
			fmt.Println(v)
		}
		fmt.Println("请输入您想下载的章节，例：第1.1话/1.1（建议直接复制）")
		var chapter string
		_, err := fmt.Scan(&chapter)
		if err != nil {
			fmt.Println("扫描出错，请稍后重试！")
			return
		}
		for _, v := range info.Groups.Default.Chapters {
			if strings.Contains(v.Name, chapter) {
				pic = info.GetChaptersInfo(v.Id)
				download.DownloadAChapter(pic, res.Result.List[index], v.Name, typeMap[v.Type])
				break
			}
		}
		//fmt.Println(pic)
	}
}
