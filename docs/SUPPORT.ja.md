# エージェントとマシンの対応状況

[简体中文](SUPPORT.zh-CN.md) · [English](SUPPORT.md) · **日本語**

「対応」には、公式パスとの互換性、インストールのライフサイクル検証、モデル動作の同等性という独立した 3 段階があります。AnyWork はこれらを分けて報告します。

| ホスト | アダプター | ターゲットファミリー | 既定で導入 | 動作同等性 |
|---|---|---|---|---|
| Codex | stable | `.agents/skills` | はい | 未検証 |
| Claude Code | stable | `.claude/skills` | はい | 未検証 |
| Gemini CLI | stable | `.agents/skills` | はい | 未検証 |
| GitHub Copilot | stable | `.agents/skills` | はい | 未検証 |
| OpenCode | stable | `.agents/skills` | はい | 未検証 |
| Cursor | experimental | `.cursor/skills` | いいえ | 未検証 |

`.agents/skills` を使う 4 ホストは、一つの物理インストールを共有します。状態ファイルは解決済みターゲットごとに利用ホストの集合を記録し、一つのホストを解除しても、別のホストが利用中のファイルは削除しません。Cursor は実環境でパスの意味、IDE と CLI の双方からの検出を確認するまで、明示指定のみとします。

標準 `SKILL.md` の frontmatter は、Agent Skills で移植可能な `name`、`description`、`license`、`compatibility`、文字列同士の `metadata` だけを使います。AnyWork はリスク情報を `allowed-tools` などの事前承認機構に変換しません。

ネイティブアーカイブは macOS、Linux、Windows の amd64/arm64 向けにビルドします。6 ターゲットへのクロスコンパイルは検証済みですが、6 種類のネイティブ runner でリリース成果物を実行する試験は未完了のリリース条件です。正式な状態は [evidence/claims.json](../evidence/claims.json)、アダプター詳細は [registry/support-matrix.json](../registry/support-matrix.json) を参照してください。
