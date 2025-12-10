# Research & Design Decisions Template

---
**Purpose**: Capture discovery findings, architectural investigations, and rationale that inform the technical design.

**Usage**:
- Log research activities and outcomes during the discovery phase.
- Document design decision trade-offs that are too detailed for `design.md`.
- Provide references and evidence for future audits or reuse.
---

## Summary
- **Feature**: `brakeman-to-codequality`
- **Discovery Scope**: New Feature (Greenfield CLI Tool)
- **Key Findings**:
  - Goの標準ライブラリ`encoding/json`の`json.Decoder`を使用した標準入力からのストリーミング処理が推奨される
  - セキュリティ上、フィンガープリント生成にはMD5ではなくSHA-256を使用すべき
  - シンプルなCLIツールであるため、CobraやViperは不要で、標準ライブラリのみで実装可能
  - テーブル駆動テストとtestifyライブラリの組み合わせが2025年のGoテストのベストプラクティス

## Research Log
Document notable investigation steps and their outcomes. Group entries by topic for readability.

### JSON処理パターン
- **Context**: 標準入力からBrakemanのJSON出力を効率的に読み込む方法の調査
- **Sources Consulted**:
  - [JSON and Go - The Go Programming Language](https://go.dev/blog/json)
  - [json package - encoding/json - Go Packages](https://pkg.go.dev/encoding/json)
  - [High-performance JSON parsing in Go](https://www.cockroachlabs.com/blog/high-performance-json-parsing/)
- **Findings**:
  - `json.NewDecoder(os.Stdin)`を使用したストリーミング処理が推奨される
  - メモリ効率の観点から、一度に全データを読み込むのではなくデコーダーを使用すべき
  - 強く型付けされた構造体を定義することでパフォーマンスとコード品質が向上
  - 不明なフィールドは無視される（`DisallowUnknownFields`を使用しない限り）
  - `json.RawMessage`を使用した部分的なパースも可能だが、今回のユースケースでは不要
- **Implications**:
  - `BrakemanReport`と`BrakemanWarning`構造体を定義して型安全性を確保
  - `json.Decoder`を使用してストリーミング処理を実装
  - エラーハンドリングでは`io.EOF`の特別な扱いが必要

### CLI設計パターン
- **Context**: GoでのCLIツール実装のベストプラクティス調査
- **Sources Consulted**:
  - [Building CLI Apps in Go with Cobra & Viper](https://www.glukhov.org/post/2025/11/go-cli-applications-with-cobra-and-viper/)
  - [GitHub - spf13/cobra: A Commander for modern Go CLI interactions](https://github.com/spf13/cobra)
  - [Building Powerful Go CLI Apps](https://golang.elitedev.in/golang/building-powerful-go-cli-apps-complete-cobra-and-viper-integration-guide-for-developers-50bcfb1e/)
- **Findings**:
  - CobraとViperは複雑なCLIアプリケーション向けのデファクトスタンダード
  - サブコマンド、フラグ、設定ファイルが必要な場合に有用
  - 今回のツールはシンプルな変換処理のみであり、標準ライブラリで十分
  - プロジェクト構造は`cmd/`ディレクトリと`internal/`ディレクトリの分離が推奨される
- **Implications**:
  - CobraやViperは使用せず、シンプルな`main.go`のみで実装
  - 将来的にフラグやサブコマンドが必要になった場合のみ、Cobraの導入を検討

### エラーハンドリングパターン
- **Context**: Goでの効果的なエラーハンドリング方法の調査
- **Sources Consulted**:
  - [The Fundamentals of Error Handling in Go](https://betterstack.com/community/guides/scaling-go/golang-errors/)
  - [Mastering Error Handling in Go](https://vickychhetri.com/2024/11/19/mastering-error-handling-in-go-best-practices-and-patterns/)
  - [Error Handling Best Practices in Go](https://blog.marcnuri.com/error-handling-best-practices-in-go)
- **Findings**:
  - Go 1.13で`errors.Is`と`errors.As`が導入され、エラーのアンラップと型マッチングが簡素化
  - Go 1.20で`errors.Join`が追加され、複数のエラーをグループ化可能
  - Rust風のResultパターンも一部で使用されているが、Go標準の明示的エラーハンドリングが主流
  - `?`オペレーターの提案は合意に至らず却下された
- **Implications**:
  - 標準的なGoのエラーハンドリング（明示的な`if err != nil`チェック）を使用
  - カスタムエラー型は不要で、標準の`error`インターフェースとフォーマット関数を使用
  - 無効な警告はスキップし、処理を継続する

### ハッシュアルゴリズム選択
- **Context**: フィンガープリント生成に使用するハッシュアルゴリズムの選択
- **Sources Consulted**:
  - [Hash Speed: MD5, SHA-1, SHA-256, SHA-512 in Go](https://blog.vitalvas.com/post/2025/10/10/hash-speed-md5-sha1-sha256-sha512-golang/)
  - [Hash Algorithm Comparison 2025](https://shautils.com/hash-algorithm-comparison)
  - [The md5 hashing algorithm is insecure](https://docs.datadoghq.com/security/code_security/static_analysis/static_analysis_rules/go-security/import-md5/)
- **Findings**:
  - MD5は2004年以降暗号学的に破られており、セキュリティ目的には使用不可
  - SHA-256は2025年時点で業界標準であり、実用的な攻撃は発見されていない
  - SHA-NI（ハードウェアアクセラレーション）によりSHA-256のパフォーマンスが2-3倍向上
  - Goの`crypto/sha256`は利用可能な場合に自動的にSHA-NIを使用
- **Implications**:
  - フィンガープリント生成にはSHA-256を使用（`crypto/sha256`パッケージ）
  - セキュリティツールの出力を扱うため、暗号学的に安全なハッシュを選択

### テストパターン
- **Context**: Goでの効果的なテスト手法の調査
- **Sources Consulted**:
  - [Go Unit Testing: Structure & Best Practices](https://www.glukhov.org/post/2025/11/unit-tests-in-go/)
  - [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests)
  - [Writing Table-Driven Tests Using Testify in Go](https://tillitsdone.com/blogs/table-driven-tests-with-testify/)
- **Findings**:
  - テーブル駆動テストがGoのテストのベストプラクティス
  - testifyライブラリの`require`パッケージは失敗時にテストを即座に停止（99%のケースで推奨）
  - `assert`は失敗後も継続するため、複数のアサーションが必要な場合のみ使用
  - 並列テスト実行には`tt := tt`のループ変数キャプチャが必要
  - `t.Run()`を使用したサブテストにより、テストの組織化と選択的実行が可能
- **Implications**:
  - 変換ロジック、マッピング関数、エラーハンドリングにテーブル駆動テストを適用
  - testifyの`require`パッケージを使用
  - JSON入出力のテストケースを複数のシナリオで網羅

## Architecture Pattern Evaluation
List candidate patterns or approaches that were considered. Use the table format where helpful.

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| 単純な手続き型 | main関数内で全処理を実装 | シンプル、依存なし | テスト困難、拡張性低い | 小規模ツールには適しているが、テストが難しい |
| 関数分離パターン | 変換ロジックを独立した関数に分離 | テスト可能、読みやすい | 軽度の複雑性増加 | **選択**：テスタビリティと保守性のバランスが最良 |
| レイヤードアーキテクチャ | 入力/変換/出力を明確に分離 | 高い保守性、明確な責務分離 | 小規模ツールには過剰 | 将来的な機能追加時に検討 |

## Design Decisions
Record major decisions that influence `design.md`. Focus on choices with significant trade-offs.

### Decision: `SHA-256をフィンガープリントに使用`
- **Context**: GitLab Code Qualityのフィンガープリントフィールドに使用するハッシュアルゴリズムの選択
- **Alternatives Considered**:
  1. MD5 — 高速だが暗号学的に安全でない
  2. SHA-256 — 業界標準、安全、十分な速度
  3. SHA-512 — より強力だが出力が長く、パフォーマンスメリットは限定的
- **Selected Approach**: SHA-256（`crypto/sha256`パッケージ）
- **Rationale**:
  - セキュリティツールの出力を扱うため、暗号学的に安全なハッシュが適切
  - SHA-256は2025年の業界標準であり、GitLabも推奨
  - ハードウェアアクセラレーション（SHA-NI）により十分な速度
- **Trade-offs**:
  - MD5より若干遅いが、セキュリティメリットが勝る
  - 出力長（64文字の16進数）は十分に短く、フィンガープリントとして適切
- **Follow-up**: パフォーマンステストで速度を検証

### Decision: `CobraやViperを使用しない`
- **Context**: CLIフレームワークの選択
- **Alternatives Considered**:
  1. Cobra + Viper — 高機能なCLIフレームワーク
  2. 標準ライブラリのみ — シンプル、依存なし
- **Selected Approach**: 標準ライブラリのみで実装
- **Rationale**:
  - ツールの要件がシンプル（標準入力→変換→標準出力）
  - サブコマンド、複雑なフラグ、設定ファイルが不要
  - 依存関係を最小化し、バイナリサイズを小さく保つ
  - 要件5.4でハッシュアルゴリズムの選択が許容されており、将来的なフラグ追加の可能性は低い
- **Trade-offs**:
  - 将来的にフラグやサブコマンドを追加する場合、リファクタリングが必要
  - 現時点の要件ではこのリスクは許容範囲内
- **Follow-up**: 将来的な機能追加要件があればCobraへの移行を検討

### Decision: `構造体による型安全な変換`
- **Context**: BrakemanとGitLab形式間のデータ変換方法
- **Alternatives Considered**:
  1. map[string]interface{}を使用した動的処理
  2. 構造体を定義した型安全な処理
- **Selected Approach**: BrakemanWarning、CodeQualityViolation構造体を定義
- **Rationale**:
  - コンパイル時の型チェックによりバグを早期発見
  - IDEの補完とリファクタリングサポート
  - ドキュメントとしての役割（構造が明確）
  - Goのベストプラクティスに準拠
- **Trade-offs**:
  - Brakemanの新しいフィールド追加時、構造体の更新が必要
  - ただし、未知フィールドは自動的に無視されるため、互換性は維持
- **Follow-up**: Brakemanのバージョンアップ時に構造体を見直し

## Risks & Mitigations
- **Brakemanのスキーマ変更** — 構造体のフィールドはオプショナルとし、未知フィールドを無視することで後方互換性を確保
- **大容量JSON入力** — ストリーミング処理（`json.Decoder`）により、メモリ使用量を抑制
- **無効な警告データ** — 要件7.3に従い、無効な警告をスキップして処理を継続することで、部分的な変換結果を提供

## References
Provide canonical links and citations (official docs, standards, ADRs, internal guidelines).
- [JSON and Go](https://go.dev/blog/json) — 公式JSON処理ガイド
- [encoding/json package](https://pkg.go.dev/encoding/json) — 標準JSONパッケージドキュメント
- [crypto/sha256 package](https://pkg.go.dev/crypto/sha256) — SHA-256ハッシュパッケージ
- [GitLab Code Quality Documentation](https://docs.gitlab.com/ci/testing/code_quality/) — GitLab Code Quality形式仕様
- [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) — テーブル駆動テストの公式ガイド
