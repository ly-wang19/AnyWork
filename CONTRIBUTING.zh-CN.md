# 参与 AnyWork

**简体中文** · [English](CONTRIBUTING.md) · [日本語](CONTRIBUTING.ja.md)

请贡献聚焦任务、来源清晰、安全且可测试的工作能力。不要提交提示词大杂烩，也不要复制许可证不兼容或仅仅公开可见的内容。

## 提交 PR 前

1. 选择一个原子工作结果，并使用小写连字符 Skill 名称。
2. 只维护一份 `SKILL.md` 工作流逻辑；在触发描述和评测中覆盖英文、简体中文、日文。
3. 在 `registry/catalog.json` 声明来源、许可证、版本、风险和副作用。
4. 改编第三方内容时，在 `THIRD_PARTY.yml` 记录不可变上游版本、作者、许可证、署名和修改说明。
5. 添加三语评测；存在确定性行为时添加自动化测试。
6. 运行：

   ```bash
   python3 anywork-cli.py doctor --lang zh-CN
   PYTHONPATH=src python3 -m unittest discover -s tests -v
   ```

7. 每个提交都加入 `Signed-off-by: 姓名 <邮箱>`，确认遵守 [DCO](DCO) 中的 Developer Certificate of Origin 1.1。

自动翻译只能作为草稿；正式版要求三种语言都经过合格人工审校。AI 辅助贡献的权利、事实、安全和测试责任仍由提交者承担。

## 第三方边界

- `original`：原创内容，按 Apache-2.0 和 DCO 贡献。
- `adapted`：仅在许可证兼容且来源完整时复制或改编。
- `curated`：只保存元数据和上游链接，不复制内容。

不得把无许可证、仅源码可见、非商业、禁止演绎或其他不兼容内容复制进核心发行物。
