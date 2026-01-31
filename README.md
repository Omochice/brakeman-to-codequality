# brakeman-to-codequality

A Go command-line tool that converts [Brakeman](https://brakemanscanner.org/) security scan results to [GitLab Code Quality](https://docs.gitlab.com/ee/ci/testing/code_quality.html) format.

## Overview

This tool reads Brakeman JSON output from standard input and outputs GitLab Code Quality JSON format to standard output, making it easy to integrate Brakeman security scans into GitLab CI/CD pipelines.

## Installation

### From Source

```bash
go install github.com/Omochice/brakeman-to-codequality@latest
```

### Build Locally

```bash
git clone https://github.com/Omochice/brakeman-to-codequality.git
cd brakeman-to-codequality
go build
```

## Usage

### Basic Usage

```bash
brakeman -f json | brakeman-to-codequality > codequality.json
```

### With File Input/Output

```bash
brakeman -f json -o brakeman-report.json
cat brakeman-report.json | brakeman-to-codequality > codequality.json
```

## CI/CD Integration

### GitLab CI Example

```yaml
brakeman:
  stage: test
  image: ruby:3.2
  before_script:
    - gem install brakeman
    - go install github.com/Omochice/brakeman-to-codequality@latest
  script:
    - brakeman -f json | brakeman-to-codequality > codequality.json
  artifacts:
    reports:
      codequality: codequality.json
```

## Format Conversion Details

### Severity Mapping

Brakeman confidence levels are mapped to GitLab severity levels:

- **High** → `critical`
- **Medium** → `major`
- **Weak** → `minor`
- **Low** → `minor`
- Unknown → `info`

### Fingerprint Generation

Each violation receives a unique SHA-256 fingerprint based on:
- File path
- Line number
- Warning type
- Message
- Code snippet (if available)

This ensures GitLab can accurately track warnings across scans.

## Exit Codes

- `0`: Success
- `1`: Error (invalid JSON, I/O error, etc.)

## Error Handling

- Invalid or missing required fields in Brakeman warnings are skipped
- Error messages are written to standard error
- Empty warning arrays produce valid empty GitLab Code Quality output

## Requirements

- Go 1.21 or later

## License

[zlib](./LICENSE)

## Contributing

Contributions are welcome! Please open an issue or pull request.
