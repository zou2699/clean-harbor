# clean-harbor

半智能清理harbor中的镜像

### Usage

```shell
Usage of clean-harbor:
  -h    help message
  -keepNum int
        每个repo保留的tag个数 (default 5)
  -password string
        harbor账号
  -projectName string
        projectName
  -url string
        harbor地址
  -user string
        harbor账号
```

### build

```shell
go build .
```

### crontab

```shell
for example
#> crontab -l
0 2 * * * /root/clean-harbor -url http://10.0.0.1 -user clean -password password -projectName cloud -keepNum 5 >> /var/log/cleanHarbor.log
```



```shell



```


