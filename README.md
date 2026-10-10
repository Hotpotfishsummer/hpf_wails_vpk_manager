# L4D2 VPK Manager (Go + Wails + Vue3 版本)

求生之路 2 (Left 4 Dead 2) VPK 文件管理工具的现代化原生重构版本，采用 **Go + Wails v2 + Vue 3 + TypeScript** 架构构建。

从原先的 Python + PyWebView + Nuitka 方案全盘迁移而来，打包产物由原来的 50~80MB 缩减至 **~14MB**，构建速度从数分钟提升至仅需 **10 多秒**，具备更高的运行速度、更低的内存占用和原生并发能力。

---

## ✨ 功能特性

* 📁 **L4D2 目录管理**：浏览并选择游戏安装根目录，自动检测 `addons` 与 `workshop`。
* ⚡ **超高速并发扫描**：利用 Go Goroutine 并发池快速解析本地及工坊 VPK，秒级加载数千个模组。
* 📦 **高性能 VPK 元数据解析**：自研纯 Go VPK v1/v2 格式解析器，支持多编码容错提取 `addoninfo.txt` 标题。
* 💾 **智能元数据缓存**：与原版本保持完全兼容的 `~/.vpk_config/<vpk_stem>.json` 校验缓存，避免重复解析。
* 🖼️ **缩略图同步提取**：自动匹配并渲染同名 `.jpg` / `.png` / `.jpeg` 图片。
* 🗜️ **多线程压缩导出**：将选中 VPK 及关联图片导出为 `.tar.zst`，采用纯 Go 极致优化的 `zstd` 算法并发压缩。
* 📥 **多格式解压导入**：支持从 `.zip`、`.7z`、`.tar.zst` 导入 VPK 和关联缩略图，自带路径穿越安全防护。
* 🗑️ **禁用与删除联动**：支持一键切换 `.vpk.disabled` 状态，自动同步重命名/删除关联图片与说明文本。
* 🌐 **SFTP 远程同步**：纯 Go 原生 SSH/SFTP 客户端，支持远程文件浏览、上传、下载、重命名与删除。
* 🔄 **L4D2 服务重启**：一键调用远程重启接口重载服务器容器。
* ⚙️ **配置模板自动化**：自动对比、备份并应用 `autoexec.cfg` 模板。
* 🌍 **中英双语国际化**：完整的 vue-i18n 多语言支持。

---

## 🛠️ 技术栈

* **后端：** Go 1.22+ (内置 Goroutines, pure Go VPK parser, klauspost/compress/zstd, pkg/sftp, bodgit/sevenzip)
* **框架：** [Wails v2](https://wails.io/) (基于原生 Windows WebView2)
* **前端：** Vue 3 + Vite + TypeScript + Pinia + Vue-i18n
* **编译：** 原生单一 EXE 可执行文件，免外部运行环境

---

## 🚀 快速上手

### 1. 开发模式 (热重载)

直接在项目根目录下运行一键开发脚本：

```powershell
.\scripts\dev.ps1
```

或者使用 Wails CLI：

```powershell
wails dev
```

### 2. 生产构建 (生成单个 .exe)

运行一键构建脚本：

```powershell
.\scripts\build.ps1
```

或者使用 Wails CLI：

```powershell
wails build -clean
```

构建生成的单文件可执行程序位于：
`build/bin/L4D2-VPK-Manager.exe`

### 3. 运行测试

```powershell
# 运行全部 Go 后端单元测试与集成测试
go test -v ./...
```
