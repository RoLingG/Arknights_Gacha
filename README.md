# Arknights_Gacha_History

明日方舟抽卡记录，代码里面该有的都有，前端是靠AI糊的，有兴趣的可以根据源代码看看相关的接口以及对应的请求参数。运行前记得将`template.env`改成`.env`，并配置好 `.env` 里面的手机号和密码

## 后端

### Windows11

仓库拉下来直接在仓库目录下的`Cmd`运行：

```cmd
## 打包程序
go build main.go

## 直接运行
.\main.exe
```

### Linux

仓库拉下来之后先 `build` 后再运行：

```bash
## 打包程序
go build -o gacha main.go

## 直接运行
./gacha

## 后台运行
nohup ./gacha &
```

## 前端

本地运行直接在浏览器 URL 输入 `localhost` 加上对应配置的端口就可以了：

```url
localhost:8081
```

对了，如果要部署的话，记得改改 `script.js` 里面获取后端的子域名 `localhost` ，改成自己的域名，不然会获取不到接口数据。

如果有什么问题的话可以直接提 `issue`，可能会改可能不会改，大概率只会改改前端（

大概就这样，谢谢你进来看我的仓库和我写的代码，希望没有污染到你的眼睛。

### 新更新

现在更新了 Vue3 前端，变化很大：

```bash
cd Gacha_Vue
npm install
# 本地测试直接进 dev 给的 localhost 看就好了
npm run dev
# 或者
npm run build
# build 完之后记得将dist目录下打包好的文件覆盖外层的文件
```

后端开启好，前端也跑好了，就可以看新的页面设计辣！！！！
