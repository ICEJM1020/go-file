# Based On: [songquanpeng/go-file](https://github.com/songquanpeng/go-file), RESPECT!

<p align="center">
  <a href="https://github.com/songquanpeng/go-file"><img src="https://user-images.githubusercontent.com/39998050/108494937-1a573e80-72e3-11eb-81c3-5545d7c2ed6e.jpg" width="200" height="200" alt="go-file"></a>
</p>

<div align="center">

# Go File

_✨ 文件分享工具，仅单个可执行文件，开箱即用，可用于局域网内分享文件和文件夹，直接跑满本地带宽 ✨_  

</div>

<p align="center">
  <a href="https://raw.githubusercontent.com/songquanpeng/go-file/master/LICENSE">
    <img src="https://img.shields.io/github/license/songquanpeng/go-file?color=brightgreen" alt="license">
  </a>
  <a href="https://github.com/songquanpeng/go-file/releases/latest">
    <img src="https://img.shields.io/github/v/release/songquanpeng/go-file?color=brightgreen&include_prereleases" alt="release">
  </a>
  <a href="https://github.com/songquanpeng/go-file/releases/latest">
    <img src="https://img.shields.io/github/downloads/songquanpeng/go-file/total?color=brightgreen&include_prereleases" alt="release">
  </a>
  <a href="https://hub.docker.com/repository/docker/justsong/go-file">
    <img src="https://img.shields.io/docker/pulls/justsong/go-file?color=brightgreen" alt="docker pull">
  </a>
  <a href="https://goreportcard.com/report/github.com/songquanpeng/go-file">
  <img src="https://goreportcard.com/badge/github.com/songquanpeng/go-file" alt="GoReportCard">
  </a>
</p>

## 特点
1. 无需配置环境，仅单个可执行文件，直接双击即可开始使用。
2. 自动打开浏览器，分享文件快人一步。
3. 提供二维码，可供移动端扫描下载文件，告别手动输入链接。
4. 支持分享本地文件夹。
5. 适配移动端。
6. 内置图床，支持直接粘贴上传图片，提供图片上传 API。
7. 内置视频播放页面，可用于在其他设备上在线博客自己电脑上的视频，轻松跨设备在线看视频。
8. 支持拖拽上传，拷贝上传。
9. 允许对不同类型的用户设置文件访问权限限制。
10. 访问频率限制。
11. 支持 Token API 验证，便于与其他系统整合。
12. 为不熟悉命令行的用户制作了启动器，[详见此处](https://github.com/songquanpeng/gofile-launcher)。
13. 配套 CLI 工具，支持命令行上传文件，支持 P2P 模式文件分享，[详见此处](https://github.com/songquanpeng/gofile-cli)。
14. Docker 一键部署：`docker run -d --restart always -p 3000:3000 -e TZ=Asia/Shanghai -v /home/ubuntu/data/go-file:/data justsong/go-file`

# TimberVersion

## New Features
1. WebDAV support, now available for Infuse. Need a account, if you want, set a 'webdav' acount for WebDAV only.
\[TO-DO\]
2. Offline download using Aria2-Pro (Thanks to https://github.com/P3TERX/Aria2-Pro-Docker)



