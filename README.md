# L4D2 VPK Manager (Go + Wails + Vue3 版本)

求生之路 2 (Left 4 Dead 2) VPK 文件管理工具的现代化原生重构版本，采用 **Go + Wails v2 + Vue 3 + TypeScript** 架构构建。

从原先的 Python + PyWebView + Nuitka 方案全盘迁移而来，打包产物由原来的 50~80MB 缩减至 **~14MB**，构建速度从数分钟提升至仅需 **10 多秒**，具备更高的运行速度、更低的内存占用和原生并发能力。

---

## ✨ 核心特性

### 📁 客户端与目录管理
* 🔍 **多源游戏客户端自动探测**：在 0.1 秒内自动扫描系统运行进程、Windows 注册表、开始菜单/桌面快捷方式及各磁盘 Steam 游戏库，并支持全盘深度探测，一键切换不同客户端目录（完美兼容 Steam 版与各类对战平台/独立免安装版）。
* 🚀 **游戏启动器与启动项优化**：原生一键启动游戏（Steam 版智能走 Steam 协议保证联机与云存档，独立版直接执行主程序），支持在设置中开启**跳过片头动画 (`-novid`)**、**强制简体中文 (`-language schinese`)**、**直接运行模式**以及自定义启动命令行参数。

### 📦 模组管理与并发性能
* ⚡ **超高速并发扫描**：利用 Go Goroutine 并发池快速解析本地及工坊 VPK，秒级加载数千个模组。
* 🏷️ **自研 VPK 元数据解析器**：自研纯 Go VPK v1/v2 格式解析器，支持多编码容错提取 `addoninfo.txt` 标题与缩略图。
* 💾 **智能元数据缓存**：基于文件校验的 `~/.vpk_config/<vpk_stem>.json` 缓存机制，避免重复解析。
* 🗂️ **VPK 自定义分组系统**：
  * **一键分组与管理**：勾选多个模组可一键创建分组、加入组、移出组或解散组；
  * **视觉隔离与色彩标识**：提供 12 款暗黑电竞调色盘，同一分组模组拥有统一高亮背景色与边框标识，并在列表中强制聚合相邻排列；
  * **折叠与统计**：支持分组一键折叠/展开，显示组内模组数量，并支持全选组内文件批量操作。
* 🤖 **OpenAI 标准大模型接入基座**：
  * 内置标准 OpenAI Chat Completions (`/v1/chat/completions`) 协议基座，通用支持 **OpenAI (GPT-4o/mini)**、**DeepSeek**、**阿里云通义千问**、**Kimi / Moonshot** 以及本地模型框架（**Ollama**、**LM Studio**、**LocalAI** 等）；
  * 设置中心支持灵活配置 Base URL、API Key、模型名称、采样温度与最大 Tokens，提供一键连接测试与延迟反馈。
* 🧠 **双轨分组名称智能自动预测**：
  * **AI 智能归纳**：配置 API Key 时，自动聚合模组文件名、标题、战役信息和备注，由大模型精准提炼精炼贴切的中文分组名称；
  * **本地领域启发式规则引擎**：离线或未配置 Key 时，秒级自动分析生还者角色（佐伊、尼克等）、特感（Tank、Witch等）、武器枪械、战役地图、音效与公共词根，0 延迟保底输出合理名称。

### 🗺️ 地图、备注与文件生态
* 🗺️ **战役章节穿透解析**：自动解析 VPK 内部包含的全部第三方地图章节代号（如 `c1m1_hotel`），提供一键复制控制台 `map` 启动指令，轻松解决三方图无主菜单海报问题。
* 📝 **自定义模组中文备注**：支持为任意本地或工坊模组保存自定义中文备注，备注自动与标题融合并支持顶部检索联动。
* 🗜️ **多线程压缩打包导出**：采用高效 `zstd` 算法将所选 VPK 及关联缩略图快速归档为 `.tar.zst`。
* 📥 **多格式解压导入 (支持拖拽)**：支持直接拖拽 `.zip`、`.7z`、`.tar.zst` 或 `.vpk` 文件至窗口任意区域自动解压导入，具备路径穿越安全防护。
* ⚡ **极速启停与安全删除**：底层通过 `.vpk.disabled` 状态切换实现瞬时启停，零耗时零磁盘磨损；删除时同步联动清理关联缩略图与配置。

### 🌐 远程服务器与运维
* 🌐 **SFTP 远程双向同步**：内置纯 Go 原生 SSH/SFTP 客户端，支持双栏分屏工作区、自由拖拽分割线、远程文件上传、下载、重命名与删除。
* 🔄 **L4D2 服务器容器重启**：一键调用远程容器管理接口，快速重启游戏服务器并实时反馈耗时。
* ⚙️ **配置模板自动化**：自动对比、备份并应用 `autoexec.cfg` 模板。
* 🌍 **中英双语国际化**：完整的 vue-i18n 多语言无缝切换支持。

---

## 🛠️ 技术栈

* **后端：** Go 1.22+ (Goroutines 并发池, 自研纯 Go VPK 解析器, klauspost/compress/zstd, pkg/sftp, bodgit/sevenzip)
* **框架：** [Wails v2](https://wails.io/) (基于原生 Windows WebView2)
* **前端：** Vue 3 + Vite + TypeScript + Pinia + Vue-i18n
* **编译：** 原生单一 EXE 可执行文件，免外部运行环境与 Python 运行时

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
`build/bin/L4D2-VPK-Manager.exe` (~14.4 MB)

### 3. 运行测试

```powershell
# 运行全部 Go 后端单元测试与集成测试
go test -v ./...
```

### 4. 自动化版本发布 (GitHub Actions CI/CD)

项目已配置全自动 Release 工作流：
* **检测版本号自动发布**：只要修改 `wails.json`（或 `frontend/package.json`）中的版本号（如 `1.0.0` -> `1.0.1`）并推送到 `main` 分支，GitHub Actions 会自动检测版本变更、启动 Windows 构建打包并创建附带 `.zip` 与 `.sha256` 的全新 GitHub Release。
* **手动触发发布**：亦可在 GitHub 仓库的 **Actions -> Build and Release** 页面点击「Run workflow」手动触发。
* **Tag 推送发布**：直接推送 Git Tag（如 `git tag v1.0.1 && git push origin v1.0.1`）同样会自动构建并发布。

---

## 📄 配置文件说明

应用程序在用户系统 `Documents/L4D2-VPK-Manager` 目录下存储用户偏好数据：
- `preferences.json`：存储上次打开路径、SFTP 配置、启动参数选项与 OpenAI 大模型 API 配置。
- `vpk_groups.json`：持久化存储模组分组、分组颜色与映射关系。
- `vpk_notes.json`：存储模组的自定义备注文本。
