# 実装タスク

## タスク概要

本ドキュメントは、brakeman-to-codequalityツールの実装タスクを定義します。各タスクは1-3時間程度で完了可能な単位に分割されており、要件と設計に基づいて段階的に実装を進めます。

## 実装タスク

- [ ] 1. プロジェクトセットアップとデータ構造定義
- [x] 1.1 Goモジュールの初期化とプロジェクト構造の作成
  - `go mod init`でモジュールを初期化
  - `main.go`を作成し、基本的なパッケージ構造を設定
  - `go.mod`に必要な依存関係（testify）を追加
  - _Requirements: 1.1_

- [x] 1.2 Brakeman入力データ構造の定義
  - `BrakemanReport`構造体を定義（warnings配列を含む）
  - `BrakemanWarning`構造体を定義（warning_type, message, file, line, confidence, codeフィールド）
  - JSONタグを適切に設定（`json:"field_name"`および`omitempty`）
  - _Requirements: 2.1, 2.2, 2.3_

- [x] 1.3 GitLab出力データ構造の定義
  - `CodeQualityViolation`構造体を定義（description, check_name, fingerprint, severity, location）
  - `Location`および`Lines`構造体を定義
  - JSONタグを適切に設定
  - _Requirements: 3.1, 3.3_

- [ ] 2. 重要度マッピング機能の実装
- [x] 2.1 (P) MapSeverity関数の実装
  - Brakemanの信頼度（High/Medium/Weak/Low）をGitLabの重要度（critical/major/minor/info）にマッピング
  - 大文字小文字を無視した比較を実装（`strings.ToLower()`）
  - 未知の信頼度レベルに対してデフォルト値"info"を返す
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [x] 2.2* (P) MapSeverity関数のユニットテスト
  - テーブル駆動テストで全信頼度レベルのマッピングを検証
  - 大文字小文字の違いをテスト（"high", "HIGH", "High"など）
  - 未知の値のデフォルト処理をテスト
  - testify/requireを使用してアサーション
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 3. フィンガープリント生成機能の実装
- [x] 3.1 (P) GenerateFingerprint関数の実装
  - ファイルパス、行番号、警告タイプ、メッセージを結合した入力文字列を生成
  - codeフィールドが存在する場合は含める
  - `crypto/sha256`を使用してSHA-256ハッシュを計算
  - `hex.EncodeToString()`で64文字の16進数文字列に変換
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 3.2* (P) GenerateFingerprint関数のユニットテスト
  - 同一入力で一貫したハッシュが生成されることを検証
  - codeフィールドありなしの両方をテスト
  - 空フィールドの処理をテスト
  - 異なる入力で異なるハッシュが生成されることを検証
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [ ] 4. 変換ロジックの実装
- [x] 4.1 ConvertWarnings関数の実装
  - BrakemanWarningのスライスを受け取り、CodeQualityViolationのスライスを返す
  - 各警告に対してMapSeverityとGenerateFingerprintを呼び出す
  - フィールドマッピングを実装（warning_type→check_name, message→description, file→location.path, line→location.lines.begin）
  - ファイルパスから"./"プレフィックスを除去
  - 必須フィールド（file, line, warning_type, message）が欠けている警告はスキップ
  - _Requirements: 3.2, 3.4, 3.5, 3.6, 3.7, 3.8, 7.3_

- [ ] 4.2* ConvertWarnings関数のユニットテスト
  - 有効な警告の正しい変換を検証
  - 無効な警告（必須フィールド欠落）のスキップを検証
  - 空配列入力の処理を検証
  - ファイルパス正規化（"./"除去）を検証
  - 複数の警告を含む配列の処理を検証
  - _Requirements: 3.2, 3.4, 3.5, 3.6, 3.7, 3.8, 7.3_

- [ ] 5. 入力処理の実装
- [x] 5.1 ParseBrakemanJSON関数の実装
  - `json.NewDecoder(io.Reader)`を使用してストリーミング処理
  - BrakemanReport構造体にデコード
  - warningsフィールドが存在しない場合は空のスライスとして扱う
  - デコードエラーを適切に返す
  - _Requirements: 1.1, 1.2, 2.1, 2.4_

- [ ] 5.2* ParseBrakemanJSON関数のユニットテスト
  - 有効なJSON構造のパースを検証
  - 無効なJSONのエラー検出を検証
  - 空入力のエラー検出を検証
  - warningsフィールドがない場合の空スライス処理を検証
  - _Requirements: 1.2, 1.3, 1.4, 2.4_

- [ ] 6. 出力処理の実装
- [x] 6.1 WriteCodeQualityJSON関数の実装
  - `json.NewEncoder(io.Writer)`を使用してJSON配列を出力
  - BOMなしUTF-8出力を保証（標準ライブラリが自動対応）
  - 空配列も正しく出力
  - エンコードエラーを適切に返す
  - _Requirements: 6.1, 6.2, 6.3_

- [ ] 6.2* WriteCodeQualityJSON関数のユニットテスト
  - 正しいJSON配列形式の出力を検証
  - 空配列の出力を検証
  - BOMなし出力を検証（bytes.Bufferで検証）
  - 複数の違反を含む配列の出力を検証
  - _Requirements: 6.1, 6.2, 6.3_

- [ ] 7. エラーハンドリングの実装
- [x] 7.1 HandleError関数の実装
  - エラーメッセージを`fmt.Fprintf(os.Stderr, ...)`で標準エラーに出力
  - 人間が読みやすい形式でエラーを整形
  - コンテキスト情報（エラーの種類、詳細）を含める
  - 非ゼロの終了コード（1）を返す
  - _Requirements: 7.1, 7.2_

- [ ] 7.2* HandleError関数のユニットテスト
  - エラーメッセージが標準エラーに出力されることを検証
  - 非ゼロの終了コードが返されることを検証
  - エラーメッセージの形式を検証
  - _Requirements: 7.1, 7.2_

- [ ] 8. メイン関数の統合
- [x] 8.1 main関数の実装
  - ParseBrakemanJSON(os.Stdin)で標準入力から読み込み
  - エラー時はHandleErrorを呼び出し、返された終了コードで終了
  - ConvertWarningsで変換を実行
  - WriteCodeQualityJSON(violations, os.Stdout)で標準出力に書き込み
  - 成功時は終了コード0で終了
  - _Requirements: 1.1, 6.4, 7.1_

- [ ] 8.2* メイン関数の統合テスト
  - 完全なエンドツーエンド変換フローを検証
  - 標準入力からBrakeman JSONを読み込み、標準出力にGitLab形式を出力
  - テストファイルを使用して実際のBrakeman出力をテスト
  - エラーシナリオ（無効なJSON、空入力）を検証
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 6.1, 6.4, 7.1, 7.4_

- [ ] 9. 追加テストとドキュメント
- [ ] 9.1* 大容量データのパフォーマンステスト
  - 数百の警告を含むJSONのメモリ効率を検証
  - 処理速度を測定（1000警告/秒以上の目標）
  - ストリーミング処理が正しく機能することを確認
  - _Requirements: すべて（パフォーマンス検証）_

- [x] 9.2 READMEの作成
  - ツールの概要と使用方法を記載
  - インストール手順を記載
  - 使用例を記載（パイプ経由での使用）
  - CI/CDパイプラインでの統合例を記載
  - _Requirements: なし（ドキュメント）_

## タスク実行ガイドライン

### コミット要件
- 各タスク完了後、必ずgitコミットを作成すること
- コミットメッセージは[Conventional Commits](https://www.conventionalcommits.org/)形式に従うこと
- コミットメッセージには変更内容とその理由を記載すること

### 並列実行可能なタスク
- タスク2.1と3.1は独立しており、並列実行可能（`(P)`マーク）
- タスク2.2と3.2も並列実行可能（`(P)`マーク）

### タスクの依存関係
- タスク1は全ての前提条件
- タスク2とタスク3はタスク1の後に実行可能
- タスク4はタスク1、2、3の後に実行
- タスク5、6、7はタスク1の後に並列実行可能
- タスク8はタスク4、5、6、7の後に実行
- タスク9はタスク8の後に実行

### テストタスクの注記
- `*`マークは、対応する実装タスクが完了した後に実行可能なオプショナルなテストカバレッジタスクを示します
- これらのテストは受入基準を検証しますが、MVP（最小viable product）のデリバリーには必須ではありません
- 実装タスクと同じ要件IDを参照しており、機能的には実装によって既に満たされています

### 推奨実行順序
1. タスク1.1, 1.2, 1.3を順次実行（データ構造定義）
2. タスク2.1と3.1を並列実行（重要度マッピングとフィンガープリント生成）
3. タスク4.1を実行（変換ロジック）
4. タスク5.1、6.1、7.1を並列実行（入力・出力・エラー処理）
5. タスク8.1を実行（メイン関数統合）
6. タスク2.2、3.2、4.2、5.2、6.2、7.2、8.2を並列実行（ユニットテストと統合テスト）
7. タスク9.1、9.2を実行（パフォーマンステストとドキュメント）
