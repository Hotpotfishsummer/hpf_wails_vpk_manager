# L4D2 Mod Manager (求生之路2 模组管理器)

<p align="center">
  <img src="frontend/src/assets/app_icon.png" width="96" height="96" alt="L4D2 Mod Manager Icon" />
</p>

<p align="center">
  <strong>专为《求生之路 2》(Left 4 Dead 2) 打造的现代化、极速、原生单文件模组与 VPK 管理器</strong>
  <br />
  A modern, ultra-fast, native single-binary Mod & VPK Manager for Left 4 Dead 2, powered by <strong>Go + Wails v2 + Vue 3 + TypeScript</strong>.
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/Hotpotfishsummer/hpf_wails_vpk_manager?color=10b981&logo=github&style=flat-square" alt="GitHub release" />
  <img src="https://img.shields.io/badge/Platform-Windows%20x64-0284c7?style=flat-square&logo=windows" alt="Platform Windows" />
  <img src="https://img.shields.io/badge/Language-Go%201.22%2B%20%7C%20TypeScript-blue?style=flat-square&logo=go" alt="Go and TS" />
  <img src="https://img.shields.io/badge/Framework-Wails%20v2-e11d48?style=flat-square" alt="Wails v2" />
  <img src="https://img.shields.io/badge/License-MIT-gray?style=flat-square" alt="License" />
</p>

---

## 📖 项目简介

**L4D2 Mod Manager**（原名 VPK Manager）是为《求生之路 2》(Left 4 Dead 2) 玩家与服务器管理员打造的专业级模组管理工具。全面采用 **Go + Wails v2 + Vue 3 + TypeScript** 架构原生构建，告别臃肿的运行环境与卡顿体验：

* 🚀 **极小体积 & 原生极速**：全编译单一 EXE 仅 **~14MB**，启动时间低于 0.3 秒，运行时内存占用仅数十兆。
* ⚡ **多协程高并发扫描**：利用 Go Goroutine 并发池快速解析本地及工坊数千个模组，秒级完成扫描与缓存。
* 🎯 **开箱即用**：零依赖运行，免安装 Python 或外部运行库。

---

## ✨ 核心特性

### 📁 客户端与启动管理
* 🔍 **多源游戏客户端自动探测**：0.1 秒自动扫描系统运行进程、Windows 注册表、快捷方式及各磁盘 Steam 游戏库，并支持全盘深度探测，一键无缝切换客户端目录（兼容 Steam 正版与各类对战平台/独立免安装版）。
* 🚀 **游戏启动器与启动项优化**：原生一键启动游戏（Steam 版智能走 Steam 协议唤起，独立版直接执行主程序），支持开启**跳过片头动画 (`-novid`)**、**强制简体中文 (`-language schinese`)**、**直接运行模式**以及自定义启动命令行参数。
* 🎮 **控制台指令与网络参数快捷复制 (v1.2.0)**：
  * **一键复制服务器连接指令**：一键将 `connect <ip>` 指令拷入剪贴板，进游戏直接按 `~` 粘贴回车即可加入服务器；
  * **100 Tick 竞技网络参数一键配置**：内置 100 Tick / 60 Tick / 30 Tick 竞技网络优化参数（`cl_updaterate`, `cl_cmdrate`, `cl_interp` 等），支持一键复制到控制台及设置中自定义。

### 📦 模组管理与并发性能
* 🔄 **即时刷新模组列表 (v1.2.0)**：
  * **无需重启应用**：新增顶部操作栏专属「刷新」按钮（带旋转加载动画）、路径输入框内快捷刷新按钮、模组列表工具栏刷新按钮；
  * **全局快捷键**：支持按 <kbd>F5</kbd> 或 <kbd>Ctrl+R</kbd> 即时重新扫描并更新模组列表；
  * **智能状态保护**：刷新后自动保留当前已勾选的模组，避免重新勾选。
* 🏷️ **自研 VPK 元数据解析器**：自研纯 Go VPK v1/v2 格式解析器，支持多编码容错提取 `addoninfo.txt` 标题、描述与缩略图。
* 💾 **智能元数据缓存**：基于文件校验的 `~/.vpk_config/<vpk_stem>.json` 缓存机制，毫秒级读取，杜绝重复解析。
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

### 🌐 远程服务器与运维 (SFTP)
* 🌐 **SFTP 远程双向同步**：内置纯 Go 原生 SSH/SFTP 客户端，支持双栏分屏工作区、自由拖拽分割线、远程文件上传、下载、重命名与删除。
* 📌 **置顶操作栏防遮挡 (v1.2.0)**：远程浏览器全选与批量操作栏（启用/禁用/删除/下载）独立固定在列表顶部，长列表滚动浏览时始终可见。
* 📊 **实时双进度条与速率反馈**：多文件并发上传时显示总体进度与单文件双进度条，实时显示传输速度与剩余时间，支持随时安全取消。
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

运行一键构建脚本（构建成功后会自动同步/覆盖更新至测试游戏目录）：

```powershell
.\scripts\build.ps1
```

> **提示**：如果只需将现有的构建产物单独同步到游戏目录，可随时运行：
> ```powershell
> .\scripts\copy_to_game.ps1
> # 或双击运行 .\scripts\copy_to_game.bat
> ```

或者使用 Wails CLI 原生构建：

```powershell
wails build -clean
```

构建生成的单文件可执行程序位于：
`build/bin/L4D2-Mod-Manager.exe` (~14.4 MB)

### 3. 运行测试

```powershell
# 运行全部 Go 后端单元测试与集成测试
go test -v ./...
```

### 4. 自动化版本发布 (GitHub Actions CI/CD)

项目已配置全自动 Release 工作流：
* **检测版本号自动发布**：只要修改 `wails.json`（或 `frontend/package.json`）中的版本号（如 `1.1.0` -> `1.2.0`）并推送到 `main` 分支，GitHub Actions 会自动检测版本变更、启动 Windows 构建打包并创建附带 `.zip` 与 `.sha256` 的全新 GitHub Release。
* **手动触发发布**：亦可在 GitHub 仓库的 **Actions -> Build and Release** 页面点击「Run workflow」手动触发。
* **Tag 推送发布**：直接推送 Git Tag（如 `git tag v1.2.0 && git push origin v1.2.0`）同样会自动构建并发布。

---

## 📄 配置文件说明

应用程序在用户系统 `Documents/L4D2-VPK-Manager`（或 `Documents/L4D2-Mod-Manager`）目录下存储用户偏好数据：
- `preferences.json`：存储上次打开路径、SFTP 配置、启动参数选项、服务器连接指令与 OpenAI 大模型 API 配置。
- `vpk_groups.json`：持久化存储模组分组、分组颜色与映射关系。
- `vpk_notes.json`：存储模组的自定义备注文本。

---

## 🤝 贡献与反馈

欢迎提交 Issue 与 Pull Request！如果您觉得这个工具有所帮助，请给仓库点个 ⭐ **Star** 支持一下！
