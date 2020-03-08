# golang 开发短地址服务

#### 使用教程
```bash
git clone https://github.com/bigfool-cn/short-link.git
```

项目根目录的short_link.sql导入的数据库，修改db.go数据库配置
```bash
go build -o app.exe
app.exe
```
默认访问地址为：http://ip:8000，ip为你本机ip, 端口可在mian.go中修改

#### 创建短地址接口
url： http://192.168.3.2:8000/api/shorten

method：post

data：{
     	"url": "http://www.baodu.com"
     }
     
#### 获取指定短地址信息接口
url： http://192.168.3.2:8000/api/info?shortlink=djme
method：get
 
#### 短地址跳转接口
url： http://192.168.3.2:8000/djme
method：get
 
     

