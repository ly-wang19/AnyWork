# AnyWork

**任何工作，任何智能体，任何机器。**

**简体中文** · [English](README.md) · [日本語](README.ja.md)

AnyWork 是采用 Apache-2.0 许可证的 AI 智能体工作能力层：包含 24 个原子技能，覆盖调研、分析、决策、写作、数据、会议、运营与持续改进，并提供原生跨平台安装器。

当前仓库是**实验性 Alpha**，不能把“世界第一”当成已经成立的宣传语。AnyWork 要把这个目标变成可检验的工程事实：每项能力、语言、智能体与机器声明，都必须指向可复现证据。

## 北极星契约

- **任何工作：**提供通用原子技能，以及可组合的工作包与岗位包。
- **任何智能体：**维护一份保守的标准技能源，并为不同宿主提供有文档依据的适配器。
- **任何机器：**提供原生二进制、确定性安装与更新、恢复、诊断和可恢复卸载。
- **三语同一标准：**中文、英文、日文遵守同一结构契约与正式版门槛。
- **默认可信：**强制记录来源、许可证、所需权限和风险等级。

Alpha 的真实状态：24 个技能全部为 `experimental`；中、英、日内容均为机器起草，尚待人工签署审校；跨智能体行为等价与完整原生机器矩阵尚未验证。详见[声明证据账本](evidence/claims.json)与[质量门槛](docs/QUALITY.zh-CN.md)。

默认文档化适配覆盖 Codex、Claude Code、Gemini CLI、GitHub Copilot 和 OpenCode。Cursor 为实验适配，不进入 `all`。静态评测契约包含 72 个案例、216 个三语提示变体，但不声称这些提示已经在各模型上完成实跑。

## 快速开始

从源码构建无第三方依赖的原生 CLI：

```bash
go build -trimpath -o anywork .
./anywork list --lang zh-CN
./anywork install essential --scope project --dry-run --lang zh-CN
./anywork install essential --scope project --lang zh-CN
./anywork doctor --lang zh-CN
```

CLI 也支持 `--lang en` 和 `--lang ja`；未指定时跟随操作系统语言。Windows 请构建 `anywork.exe`。标签发布流程会为 macOS、Linux、Windows 的 amd64/arm64 生成带校验和的压缩包。`anywork-cli.py` 保留为 Python 标准库参考实现与测试基准。

## 仓库结构

```text
skills/       聚焦单一任务的标准 Skills
registry/     工作包，以及来源、语言审校、能力、风险和版本元数据
adapters/     不同智能体的安装规则
evals/        多语言行为与质量评测
evidence/     机器可读的声明状态与证据
*.go          无第三方依赖的原生 AnyWork CLI
src/          Python 参考实现
schemas/      机器可读的契约
scripts/      可复现发布构建
tests/        注册表、证据、三语和安装器回归测试
```

## 当前状态

AnyWork 尚未正式发布。目录结构检查、本地安装生命周期测试和六目标交叉编译已经通过；母语人工审校、完整模型实跑评测、每个系统与架构上的原生执行、GitHub 公共远端和签名发布仍是发布门槛，不能被描述成已完成。

## 许可证

项目采用 Apache-2.0。第三方目录条目保留其原始许可证与署名；许可证不兼容的项目只提供链接，不复制内容。
