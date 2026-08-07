# AnyWork

**任何工作，任何智能体，任何机器。**

**简体中文** · [English](README.md) · [日本語](README.ja.md)

AnyWork 是面向 AI 智能体的开源工作能力层。它把经过验证的工作流、参考资料、脚本、模板和评测打包起来，让任何人都能在几分钟内为新电脑恢复一套高质量 AI 工作环境。

AnyWork 不是提示词大杂烩。官方 Skill 必须聚焦真实任务，原生支持多语言，经过评测，来源可追溯，并且能够安全安装。

## 产品承诺

- **任何工作：**提供通用工作能力，以及可组合的岗位包和行业包。
- **任何智能体：**维护一份标准技能源，首批适配 Codex 与 Claude Code。
- **任何机器：**提供确定性的安装、更新、诊断、回滚和卸载流程。
- **三语同等支持：**文档、触发描述、示例、CLI 消息和评测原生支持中文、英文、日文。
- **默认可信：**强制记录来源、许可证、所需权限和风险等级。

## 快速开始

```bash
python3 anywork-cli.py list --lang zh-CN
python3 anywork-cli.py install essential --agent all --dry-run --lang zh-CN
python3 anywork-cli.py install essential --agent all --lang zh-CN
python3 anywork-cli.py doctor --lang zh-CN
```

CLI 也支持 `--lang en` 和 `--lang ja`；未指定时自动跟随操作系统语言。
Windows 请使用 `py anywork-cli.py ...`。仓库 CLI 不依赖任何第三方运行库。

## 仓库结构

```text
skills/       聚焦单一任务的标准 Skills
packs/        可组合的通用、岗位和行业工作包
registry/     来源、许可证、语言、风险和版本元数据
adapters/     不同智能体的安装规则
evals/        多语言行为与质量评测
src/          跨平台 AnyWork CLI
schemas/      机器可读的契约
tests/        安装器与注册表回归测试
```

## 当前状态

AnyWork 正处于基础建设阶段，首批支持 Codex 与 Claude Code。公开质量标准和评测集是产品本身的一部分，而不是发布后的补充文档。

## 许可证

项目采用 Apache-2.0。第三方目录条目保留其原始许可证与署名；许可证不兼容的项目只提供链接，不复制内容。
