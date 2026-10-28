<div align="center">
  <p><a href="README.md">English</a> · <strong>简体中文</strong></p>

  <img src="assets/branding/gdenticon-logo.png" alt="gdenticon" width="900">

  <p><strong>为你的项目生成专属小头像。</strong></p>

  <p>
    <a href="https://github.com/hourinX/gdenticon/actions/workflows/test.yml"><img src="https://github.com/hourinX/gdenticon/actions/workflows/test.yml/badge.svg" alt="单元测试"></a>
    <a href="https://codecov.io/gh/hourinX/gdenticon"><img src="https://codecov.io/gh/hourinX/gdenticon/graph/badge.svg" alt="代码覆盖率"></a>
    <a href="https://pkg.go.dev/github.com/hourinX/gdenticon/v2/pkg/gdenticon"><img src="https://pkg.go.dev/badge/github.com/hourinX/gdenticon/v2/pkg/gdenticon.svg" alt="Go 文档"></a>
    <a href="https://github.com/hourinX/gdenticon/releases"><img src="https://img.shields.io/github/v/release/hourinX/gdenticon" alt="最新版本"></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/hourinX/gdenticon" alt="许可证"></a>
  </p>
</div>

`gdenticon` 可以将用户名、电子邮箱、账户ID或任意字符串转换成稳定的网格头像。相同的标识与配置始终生成相同的图像，适用于默认用户头像、开发者工具、仪表盘和测试数据。

项目同时提供 Go 软件包和命令行程序，支持镜像及平衡图案、根据哈希生成颜色或指定固定颜色、自定义网格与输出尺寸，以及 PNG、SVG、JPEG、GIF、BMP 和 TIFF 格式。

## 安装

安装命令行程序：

```sh
go install github.com/hourinX/gdenticon/v2/cmd/gdenticon@latest
```

添加到 Go 模块：

```sh
go get github.com/hourinX/gdenticon/v2@latest
```

也可以从 [GitHub Releases](https://github.com/hourinX/gdenticon/releases) 下载预编译的可执行文件。

gdenticon 需要 Go 1.26 或更高版本。使用 `gdenticon --help` 查看命令行选项；公开 Go API 请参阅 [pkg.go.dev](https://pkg.go.dev/github.com/hourinX/gdenticon/v2/pkg/gdenticon)。

## 命令行

在普通命令行模式下直接传入标识：

```sh
gdenticon alice@example.com --format png --size 512
```

不带参数运行 `gdenticon`，或者执行 `gdenticon shell`，即可进入交互命令模式。在交互模式中使用 `generate <id> [flags]` 生成头像，输入 `help` 查看命令，输入 `exit` 退出。

未指定输出目录时，程序会把生成的文件保存到操作系统临时目录下的 `gdenticon` 文件夹。Linux 通常为 `/tmp/gdenticon`，Windows 通常为 `%LOCALAPPDATA%\Temp\gdenticon`。

## 吉祥物

<p align="center">
  <img src="assets/branding/gdenticon-icon.png" alt="正在绘制镜像网格头像的 gdenticon 设计师 Gopher" width="360">
</p>

gdenticon 吉祥物是基于 Go Gopher 创作的衍生形象。原版 Go Gopher 由 [Renee French](https://reneefrench.blogspot.com/) 设计，并采用 [知识共享署名 4.0 国际许可协议](https://creativecommons.org/licenses/by/4.0/) 授权。本项目将其改编为正在绘制镜像 Identicon 网格的设计师形象。

## 许可证

源代码采用 [MIT License](LICENSE) 发布。
