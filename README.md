# brakeman-to-codequality

A Go command-line tool that converts [Brakeman](https://brakemanscanner.org/) security scan results to [GitLab Code Quality](https://docs.gitlab.com/ee/ci/testing/code_quality.html) format.

## Overview

This tool reads Brakeman JSON output from a file or stdin and outputs GitLab Code Quality JSON format to standard output, making it easy to integrate Brakeman security scans into GitLab CI/CD pipelines.

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

The tool requires exactly one argument: a file path or `-` for stdin.

### Read from File

```bash
brakeman -f json -o brakeman-report.json
brakeman-to-codequality brakeman-report.json > codequality.json
```

### Read from stdin

```bash
brakeman -f json | brakeman-to-codequality - > codequality.json
```

## CI/CD Integration

### GitLab CI Example

```yaml
brakeman:
  stage: test
  image: ruby:3.2
  before_script:
    - gem install brakeman
  script:
    - brakeman -f json -o brakeman-report.json
  artifacts:
    paths:
      - brakeman-report.json
    when: always

codequality:
  stage: test
  image: ghcr.io/omochice/brakeman-to-codequality:latest
  needs:
    - job: brakeman
      artifacts: true
  when: always
  script:
    - brakeman-to-codequality brakeman-report.json > codequality.json
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
