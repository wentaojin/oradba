# oradba
Oracle DBA 常用基本 SQL 查询命令行工具 oraz，支持 Oracle 会话、阻塞、备份、ASH 、表等信息查看

使用事项

```
把 oracle client instantclient-basic-linux.x64-19.8.0.0.0dbru.zip 上传解压，并配置 LD_LIBRARY_PATH 环境变量或者 任意可正常访问的 Oracle 环境
使用方法:

1、上传解压 instantclient-basic-linux.x64-19.8.0.0.0dbru.zip（该压缩包位于 github.com/wentaojin/transferdb Repo 中的 client 目录） 到指定目录，比如：/data1/soft/client/instantclient_19_8

2、查看环境变量 LD_LIBRARY_PATH
export LD_LIBRARY_PATH=/data1/soft/client/instantclient_19_8
echo $LD_LIBRARY_PATH

3、程序运行，配置文件示例见 conf 目录
./oradba --config config.toml
```

