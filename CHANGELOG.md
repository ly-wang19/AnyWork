# Changelog / 更新日志 / 変更履歴

## 0.3.0-alpha.1 — 2026-08-08

### English

- Added version-pinned, checksum-verifying one-command installers for macOS, Linux, and Windows, plus `anywork setup` for the complete capability pack.
- Added `orchestrate-work`, an end-to-end capability that selects, executes, checks, and hands off the smallest useful chain of atomic Skills.
- Added three trilingual ready-to-paste demos and an exact-prompt before/after benchmark runner with an evidence-based human review sheet.
- Expanded the static contract to 75 cases and 225 prompt variants, updated all three capability maps, and made documented pack counts match inherited installation results.

### 简体中文

- 加入适用于 macOS、Linux 和 Windows 的版本锁定、校验和验证一键安装器，以及安装完整能力包的 `anywork setup`。
- 加入 `orchestrate-work`：为复杂成果选择、执行、检查并交接最小必要原子 Skill 链。
- 加入 3 个三语可直接粘贴的 Demo，以及同提示安装前后评测器和证据化人工审阅表。
- 把静态契约扩展到 75 个案例、225 个提示变体，更新三语能力地图，并使文档工作包数量与继承后真实安装结果一致。

### 日本語

- macOS、Linux、Windows 向けのバージョン固定・チェックサム検証付き 1 コマンドインストーラーと、完全能力パック用の `anywork setup` を追加。
- 複雑な成果に必要最小限の個別 Skill チェーンを選択、実行、検査、引き継ぐ `orchestrate-work` を追加。
- 3 つの貼り付け可能な 3 言語 Demo と、同一プロンプトの導入前後評価ランナー、根拠付き人手レビューシートを追加。
- 静的契約を 75 ケース、225 プロンプトに拡張し、3 言語の能力マップを更新、パック数が継承後の実際の導入結果と一致するように修正。

## 0.2.0-alpha.1 — 2026-08-08

### English

- Added a prominent capability overview to every README and a complete trilingual map of all 24 Skills, 7 work packs, outputs, and composable workflows.
- Added CI coverage that prevents Skills or packs from disappearing from any language's capability map.
- Added an evidence-gated roadmap in English, Simplified Chinese, and Japanese, linked from every README.
- Added dependency-free Markdown structure and local-link validation to CI.
- Added a dependency-free native Go CLI with crash-consistent journals, automatic rollback, explicit recovery, shared-target ownership, and safe Windows state replacement.
- Added documented adapters for Codex, Claude Code, Gemini CLI, GitHub Copilot, and OpenCode; Cursor remains experimental and opt-in.
- Expanded the static suite to 72 cases and 216 English, Simplified Chinese, and Japanese prompt variants across basic, ambiguous, and conflicting difficulty levels.
- Reworked all 24 trigger descriptions with equivalent three-language use and exclusion boundaries; language content remains machine-drafted until signed human review.
- Added six-target archives, checksums, license payloads, provenance attestation, and native release-smoke gates.
- Added a machine-readable claims ledger so unverified native-language, cross-agent, native-platform, and stable-release claims stay explicit.

### 简体中文

- 在每份 README 首屏加入能力总览，并补充完整三语能力地图，展示 24 个 Skill、7 个工作包、具体产物和组合工作流。
- 加入 CI 校验，防止任何 Skill 或工作包从某个语言版本的能力地图中遗漏。
- 加入中、英、日三语证据门槛路线图，并从每份 README 提供入口。
- 在 CI 中加入无第三方依赖的 Markdown 结构与本地链接检查。
- 加入无第三方依赖的 Go 原生 CLI，支持崩溃一致事务日志、自动回滚、显式恢复、共享目标所有权和 Windows 安全状态替换。
- 加入 Codex、Claude Code、Gemini CLI、GitHub Copilot、OpenCode 的文档化适配；Cursor 保持实验性并仅允许显式选择。
- 静态评测扩充到 72 个案例、216 个中英日提示变体，覆盖基础、歧义和冲突三种难度。
- 重写全部 24 个 Skill 的三语适用与不适用边界；在人工签署审校前仍明确标记为机器起草。
- 加入六目标压缩包、校验和、许可证材料、构建来源证明和原生发布实跑门槛。
- 加入机器可读的声明证据账本，明确标记尚未验证的母语质量、跨智能体等价、原生平台实跑和正式版声明。

### 日本語

- 各 README の前半に能力概要を追加し、24 Skill、7 業務パック、成果物、組み合わせ可能なワークフローを示す完全な 3 言語能力マップを追加。
- いずれかの言語の能力マップから Skill やパックが欠落することを防ぐ CI 検査を追加。
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
