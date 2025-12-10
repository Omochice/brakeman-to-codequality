# Requirements Document

## Project Description (Input)
rubyの brakemanが出力するjsonをgitlabのcodequalityの形式に変換するツールをgoでかいて。なお入力はパイプで渡されるものとする

## Introduction
本ツールは、RubyのセキュリティスキャナーであるBrakemanが出力するJSON形式のレポートを、GitLab Code Quality形式のJSONに変換するコマンドラインツールです。標準入力からBrakemanのJSON出力を受け取り、GitLab Code Quality形式に変換したJSONを標準出力に出力します。

## Requirements

### Requirement 1: 標準入力からのJSON読み込み
**目的:** 開発者として、パイプを使用してBrakemanの出力を直接変換ツールに渡せるようにすることで、シェルスクリプトやCIパイプラインでの統合を容易にする

#### 受入基準
1.1. When the converter tool is executed, the converter shall read JSON data from standard input
1.2. When the standard input contains valid Brakeman JSON format, the converter shall parse it successfully
1.3. If the standard input contains invalid JSON, then the converter shall exit with a non-zero status code and output an error message to standard error
1.4. If the standard input is empty, then the converter shall exit with a non-zero status code and output an error message to standard error

### Requirement 2: Brakeman JSON形式の解析
**目的:** 開発者として、Brakemanの警告情報を正確に抽出できるようにすることで、セキュリティ問題を漏れなくGitLabに報告する

#### 受入基準
2.1. The converter shall extract warning information from the Brakeman JSON "warnings" array
2.2. The converter shall extract the following fields from each warning: warning_type, message, file, line, confidence
2.3. When a warning contains a "code" field, the converter shall include it in the fingerprint calculation
2.4. If the Brakeman JSON does not contain a "warnings" array, then the converter shall treat it as an empty warning list and output an empty GitLab Code Quality array

### Requirement 3: GitLab Code Quality形式への変換
**目的:** 開発者として、GitLabが認識できる形式で出力を得ることで、GitLab UIでコード品質レポートを表示できるようにする

#### 受入基準
3.1. The converter shall output a JSON array conforming to GitLab Code Quality format
3.2. The converter shall map each Brakeman warning to a GitLab Code Quality violation object
3.3. The converter shall populate the following required fields for each violation: description, check_name, fingerprint, severity, location.path, location.lines.begin
3.4. The converter shall use relative file paths without "./" prefix in the location.path field
3.5. The converter shall map Brakeman warning_type to the check_name field
3.6. The converter shall map Brakeman message to the description field
3.7. The converter shall map Brakeman file to the location.path field
3.8. The converter shall map Brakeman line to the location.lines.begin field

### Requirement 4: 重要度マッピング
**目的:** 開発者として、Brakemanの信頼度レベルをGitLabの重要度に適切に変換することで、問題の優先順位付けを正確に行う

#### 受入基準
4.1. The converter shall map Brakeman confidence levels to GitLab severity values
4.2. When Brakeman confidence is "High", the converter shall set severity to "critical"
4.3. When Brakeman confidence is "Medium", the converter shall set severity to "major"
4.4. When Brakeman confidence is "Weak" or "Low", the converter shall set severity to "minor"
4.5. If Brakeman confidence field is missing or unknown, the converter shall default to severity "info"

### Requirement 5: フィンガープリント生成
**目的:** 開発者として、同一の警告を一意に識別できるようにすることで、GitLabが警告の重複や修正を正しく追跡できるようにする

#### 受入基準
5.1. The converter shall generate a unique fingerprint for each violation
5.2. The converter shall calculate the fingerprint using a hash of: file path, line number, warning type, and message
5.3. When Brakeman warning contains a "code" field, the converter shall include it in the fingerprint calculation
5.4. The converter shall use a consistent hash algorithm (such as MD5 or SHA256) for all fingerprints

### Requirement 6: 標準出力へのJSON出力
**目的:** 開発者として、変換結果をパイプで他のツールに渡せるようにすることで、CI/CDワークフローでの柔軟な利用を可能にする

#### 受入基準
6.1. When conversion is successful, the converter shall write the GitLab Code Quality JSON array to standard output
6.2. The converter shall output valid JSON without byte order mark (BOM)
6.3. The converter shall format the JSON output as a single array at the root level
6.4. When conversion is successful, the converter shall exit with status code 0

### Requirement 7: エラーハンドリング
**目的:** 開発者として、変換エラー時に適切なフィードバックを得ることで、問題を迅速に特定し修正できるようにする

#### 受入基準
7.1. If an error occurs during conversion, then the converter shall exit with a non-zero status code
7.2. If an error occurs during conversion, then the converter shall output a human-readable error message to standard error
7.3. When a Brakeman warning contains invalid or missing required fields, the converter shall skip that warning and continue processing remaining warnings
7.4. When all Brakeman warnings are skipped due to invalid data, the converter shall output an empty GitLab Code Quality array and exit with status code 0
