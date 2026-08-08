# Changelog / 更新日志 / 変更履歴

## 0.2.0-alpha.1 — 2026-08-08

### English

- Added an evidence-gated roadmap in English, Simplified Chinese, and Japanese, linked from every README.
- Added dependency-free Markdown structure and local-link validation to CI.
- Added a dependency-free native Go CLI with crash-consistent journals, automatic rollback, explicit recovery, shared-target ownership, and safe Windows state replacement.
- Added documented adapters for Codex, Claude Code, Gemini CLI, GitHub Copilot, and OpenCode; Cursor remains experimental and opt-in.
- Expanded the static suite to 72 cases and 216 English, Simplified Chinese, and Japanese prompt variants across basic, ambiguous, and conflicting difficulty levels.
- Reworked all 24 trigger descriptions with equivalent three-language use and exclusion boundaries; language content remains machine-drafted until signed human review.
- Added six-target archives, checksums, license payloads, provenance attestation, and native release-smoke gates.
- Added a machine-readable claims ledger so unverified native-language, cross-agent, native-platform, and stable-release claims stay explicit.

### 简体中文

- 加入中、英、日三语证据门槛路线图，并从每份 README 提供入口。
- 在 CI 中加入无第三方依赖的 Markdown 结构与本地链接检查。
- 加入无第三方依赖的 Go 原生 CLI，支持崩溃一致事务日志、自动回滚、显式恢复、共享目标所有权和 Windows 安全状态替换。
- 加入 Codex、Claude Code、Gemini CLI、GitHub Copilot、OpenCode 的文档化适配；Cursor 保持实验性并仅允许显式选择。
- 静态评测扩充到 72 个案例、216 个中英日提示变体，覆盖基础、歧义和冲突三种难度。
- 重写全部 24 个 Skill 的三语适用与不适用边界；在人工签署审校前仍明确标记为机器起草。
- 加入六目标压缩包、校验和、许可证材料、构建来源证明和原生发布实跑门槛。
- 加入机器可读的声明证据账本，明确标记尚未验证的母语质量、跨智能体等价、原生平台实跑和正式版声明。

### 日本語

- 英語・簡体字中国語・日本語の証拠ゲート付きロードマップを追加し、各 README から参照可能にした。
- 第三者依存のない Markdown 構造・ローカルリンク検査を CI に追加。
- 第三者依存のない Go ネイティブ CLI に、クラッシュ整合性のあるジャーナル、自動ロールバック、明示的な復旧、共有ターゲットの所有権管理、Windows の安全な状態置換を追加。
- Codex、Claude Code、Gemini CLI、GitHub Copilot、OpenCode の文書化済みアダプターを追加。Cursor は実験扱いの明示指定のみ。
- 静的評価を 72 ケース、英語・簡体字中国語・日本語の 216 プロンプトに拡張し、基本・曖昧・競合の 3 難度を収録。
- 24 スキルすべての適用場面と対象外条件を 3 言語で同等に再構成。署名付き人手レビューまでは機械生成の初稿として明示。
- 6 ターゲットのアーカイブ、チェックサム、ライセンス資料、来歴証明、ネイティブなリリース実行ゲートを追加。
- 母語品質、エージェント間同等性、全ネイティブ環境、正式版について、未検証の主張を明示する機械可読台帳を追加。

## 0.1.0-alpha — 2026-08-07

### English

- Established the independent AnyWork repository and Apache-2.0 licensing.
- Added 24 atomic work Skills and 7 composable work packs.
- Added structural English, Simplified Chinese, and Japanese coverage for triggers, CLI output, help, documentation, and evaluation prompts.
- Added Codex and Claude Code project/user adapters with dry-run, conflict protection, idempotency, backup, update, and recoverable uninstall.
- Added registry validation, risky-pattern checks, official Skill validation, lifecycle tests, and a three-OS CI matrix.

### 简体中文

- 建立独立 AnyWork 仓库并采用 Apache-2.0。
- 加入 24 个原子工作 Skills 和 7 个可组合工作包。
- 建立英文、简体中文、日文触发、CLI 输出、帮助、文档和评测提示的结构覆盖。
- 加入 Codex 与 Claude Code 的用户级/项目级适配，以及预演、冲突保护、幂等、备份、更新和可恢复卸载。
- 加入注册表验证、危险模式检查、官方 Skill 验证、生命周期测试和三操作系统 CI 矩阵。

### 日本語

- 独立した AnyWork リポジトリを作成し、Apache-2.0 を採用。
- 24 個の個別業務スキルと 7 個の組み合わせ可能なパックを追加。
- 英語、簡体字中国語、日本語のトリガー、CLI 出力、ヘルプ、文書、評価プロンプトを追加。
- Codex と Claude Code のユーザー・プロジェクト向けアダプターに、ドライラン、競合保護、冪等性、バックアップ、更新、復元可能なアンインストールを追加。
- レジストリ検証、危険パターン検査、公式 Skill 検証、ライフサイクルテスト、3 OS の CI マトリクスを追加。
