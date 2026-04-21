# brakeman-to-codequality

A Go command-line tool that converts [Brakeman](https://brakemanscanner.org/) security scan results to [GitLab Code Quality](https://docs.gitlab.com/ee/ci/testing/code_quality.html) format.

## Overview

This tool reads Brakeman JSON output from a file or stdin and outputs GitLab Code Quality JSON format to stdout, making it easy to integrate Brakeman security scans into GitLab CI/CD pipelines.

## Installation

### Go install

```bash
go install github.com/Omochice/brakeman-to-codequality@latest
```

### Docker

```bash
docker pull ghcr.io/omochice/brakeman-to-codequality:latest
```

### Build from source

```bash
git clone https://github.com/Omochice/brakeman-to-codequality.git
cd brakeman-to-codequality
go build
```

## Usage

The tool requires exactly one argument: a file path or `-` for stdin.

```bash
# From a file
brakeman-to-codequality brakeman-report.json > codequality.json

# From stdin
brakeman -f json | brakeman-to-codequality - > codequality.json
```

### Options

| Flag | Description |
|------|-------------|
| `-v`, `--version` | Show application version |

### Exit codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | Error (invalid JSON, I/O error, etc.) |

## CI/CD Integration

### GitLab CI example

The following example runs Brakeman and converts its output to GitLab Code Quality format in a single pipeline.

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

### Severity mapping

Brakeman confidence levels are mapped to GitLab Code Quality severity levels as follows.

| Brakeman confidence | GitLab severity |
|---------------------|-----------------|
| High | `critical` |
| Medium | `major` |
| Weak | `minor` |
| Low | `minor` |
| (unknown) | `info` |

### Fingerprint

Each violation's fingerprint is taken directly from Brakeman's output. This ensures GitLab can accurately track warnings across scans.

### Skipped warnings

Warnings that lack any of the following required fields are silently skipped.

- `file`
- `line`
- `warning_type`
- `message`
- `fingerprint`

Error messages are written to stderr.

## Requirements

- Go 1.21 or later

## License

[zlib](./LICENSE)

## Contributing

Contributions are welcome. Please open an issue or pull request.
