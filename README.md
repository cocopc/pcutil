#### 工具类服务
##### hscp 服务
* 实现快速从本地拷贝文件，跳过跳板机，到服务器。实现方式为再跳板机部署http服务执行本地scp命令，所以要求跳板机到各服务器scp可用
* 例子：curl  host:port/api/htcp\?dsthostdir=$dsthost:/tmp/  -F file=@"$cpfile" 