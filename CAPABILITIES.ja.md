# AnyWork 能力マップ

[简体中文](CAPABILITIES.zh-CN.md) · [English](CAPABILITIES.md) · **日本語**

AnyWork は、曖昧な仕事を検証済みでそのまま使える成果物へ変える、再利用可能な作業方法を AI エージェントに提供します。現在の Alpha 版には **24 の個別 Skill、1 つのエンドツーエンド・オーケストレーター、7 つの導入可能な業務パック**があります。個別 Skill は単独で使え、オーケストレーターは複雑な成果に必要最小限のチェーンを選択・実行します。

## 能力一覧

| 業務領域 | 得られる成果 | 含まれる Skill |
|---|---|---|
| 定義と計画 | 検証可能な作業定義、実行計画、優先順位付き作業一覧 | [成果を明確化](skills/clarify-outcome/SKILL.md)、[作業を計画](skills/plan-work/SKILL.md)、[作業を優先順位付け](skills/prioritize-work/SKILL.md) |
| 調査と根拠 | 追跡可能な情報源一覧、構造化抽出、主張の検証結果、根拠の統合 | [情報源を収集](skills/gather-sources/SKILL.md)、[構造化情報を抽出](skills/extract-structured-info/SKILL.md)、[主張を検証](skills/verify-claims/SKILL.md)、[情報源を統合](skills/synthesize-sources/SKILL.md) |
| 分析と意思決定 | 再現可能な分析、根本原因報告、透明なシナリオ、選択肢比較、意思決定案 | [データを分析](skills/analyze-data/SKILL.md)、[根本原因を診断](skills/diagnose-root-cause/SKILL.md)、[シナリオをモデル化](skills/model-scenarios/SKILL.md)、[選択肢を比較](skills/compare-options/SKILL.md)、[意思決定を提案](skills/make-decision/SKILL.md) |
| 成果物とコミュニケーション | 意思決定ブリーフ、長文書、監査可能なワークブック、プレゼン、送信可能なメッセージ、ローカライズ済みコンテンツ | [ブリーフを作成](skills/write-brief/SKILL.md)、[文書を起草](skills/draft-document/SKILL.md)、[スプレッドシートを作成](skills/build-spreadsheet/SKILL.md)、[プレゼンを設計](skills/design-presentation/SKILL.md)、[メッセージを作成](skills/write-message/SKILL.md)、[コンテンツをローカライズ](skills/localize-content/SKILL.md) |
| 会議と運用 | 会議設計、意思決定・アクション一覧、実行可能な SOP、安全な自動化計画 | [会議を準備](skills/prepare-meeting/SKILL.md)、[決定とアクションを整理](skills/capture-decisions-actions/SKILL.md)、[SOP を作成](skills/create-sop/SKILL.md)、[反復業務を自動化](skills/automate-routine/SKILL.md) |
| 品質と改善 | 成果物レビュー、または担当者付き改善実験を含む根拠ベースの振り返り | [成果物をレビュー](skills/review-deliverable/SKILL.md)、[振り返りを実施](skills/run-retrospective/SKILL.md) |
| エンドツーエンド編成 | 明確な引き継ぎ、根拠、不確実性、権限ゲートを持つ検査済み複数段階成果 | [業務を編成](skills/orchestrate-work/SKILL.md) |

## 業務パック

業務パックは、よくある仕事に合わせた導入用の組み合わせです。一つの個別能力が複数の職種を支えるため、Skill は意図的に重複します。

| パック ID | 対象業務 | Skill 数 | 導入例 |
|---|---|---:|---|
| `essential` | 日常の調査、執筆、連携、レビュー、編成 | 10 | `anywork install essential` |
| `manager-leadership` | 基本パックに意思決定、優先順位、シナリオ、会議、振り返りを追加 | 18 | `anywork install manager-leadership` |
| `product-operations` | 基本パックに製品探索、優先順位、提供体制、継続的改善を追加 | 19 | `anywork install product-operations` |
| `research-consulting` | 基本パックに根拠重視の分析、提案、報告書、プレゼンを追加 | 17 | `anywork install research-consulting` |
| `go-to-market` | 基本パックに市場根拠、顧客判断、コンテンツ、コミュニケーション、実行を追加 | 21 | `anywork install go-to-market` |
| `people-recruiting` | 基本パックに構造化された公平で監査可能な人事・採用業務を追加 | 20 | `anywork install people-recruiting` |
| `complete` | 現行リリースの公式個別 Skill とオーケストレーターすべて | 25 | `anywork install complete` |

`--agent codex`、`--agent claude-code`、`--agent gemini-cli`、`--agent github-copilot`、`--agent opencode` で対象ホストを指定できます。`--scope project` または `--scope user` で導入範囲を選び、書き込み前には `--dry-run` を利用できます。

`anywork setup` は、文書化済みホストすべてに `complete` を 1 回で導入します。`anywork demo research`、`anywork demo meeting`、`anywork demo automation` は、そのまま貼り付けられる 3 言語の編成ジョブを表示します。

## エンドツーエンドの例

### 調査から提案まで

`orchestrate-work` → 必要なら `clarify-outcome` → `gather-sources` → `verify-claims` → `synthesize-sources` → `compare-options` → `make-decision` → `write-brief` → `review-deliverable`

成果：情報源、不確実性、代替案、判断理由を追跡できる、意思決定向けの提案。

### データから経営層向けプレゼンまで

`orchestrate-work` → `extract-structured-info` → `analyze-data` → `model-scenarios` → `build-spreadsheet` → `design-presentation` → `review-deliverable`

成果：根拠と仮定を分離した、再計算可能なモデルと対象者向けのストーリー。

### 会議から責任ある実行まで

`orchestrate-work` → `prepare-meeting` → 会議実施 → `capture-decisions-actions` → `plan-work` → `prioritize-work` → `review-deliverable`

成果：焦点を絞った会議と、それに続く明確な決定、担当、期限、依存関係、受入確認。

### 反復運用から安全な自動化まで

`orchestrate-work` → `diagnose-root-cause` → `create-sop` → `automate-routine` → 統制された実行 → `run-retrospective`

成果：例外、人の承認、復旧、改善測定を備えた、統制された運用プロセス。

### グローバルコンテンツの提供

`orchestrate-work` → `write-brief` → `draft-document` → `localize-content` → `write-message` → `review-deliverable`

成果：承認済みの意味と用語を維持した、英語・簡体字中国語・日本語のチャネル向けコンテンツ。

## 能力の境界

AnyWork が導入するのは作業方法であり、認証情報や暗黙の連携ではありません。Skill が使えるのは、ホストエージェントがすでに利用できるツールとアクセス権だけです。外部書き込み、メッセージ送信、購入、削除など影響を伴う操作には、引き続き明示的な承認と検証可能な対象が必要です。[品質ゲート](docs/QUALITY.ja.md)と[ロードマップ](ROADMAP.ja.md)の証拠が揃うまで、24 の個別 Skill とオーケストレーターはすべて `experimental` のままです。
