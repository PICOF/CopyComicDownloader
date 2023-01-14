# CopyComicDownloader
基于go语言的一个拷贝漫画下载器

## 使用方法
命令行参数输入 `-h` 查看全部参数
```shell
C:\Users\ANGER\goLandCode\sese> ./sese.exe -h           
Usage of C:\Users\ANGER\goLandCode\sese\sese.exe:
  -all
        是否整本下载，开启后无法指定章节
  -concurrent
        是否在整本下载时针对每个章节都采用并发下载（速度快，但容易因为请求过多导致问题）
  -limit int
        搜索漫画的返回结果条数,默认为十条 (default 10)
  -name string
        漫画名字，默认为电锯人 (default "电锯人")
```
其中 `all` 和 `concurrent` 都默认关闭