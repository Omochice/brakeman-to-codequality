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

### Docker

Build the Docker image locally:

```bash
docker build -t brakeman-to-codequality:latest .
```

Or pull from GitLab Container Registry (if using GitLab CI):

```bash
docker pull $CI_REGISTRY_IMAGE:latest
```

Run using Docker:

```bash
brakeman -f json | docker run --rm -i brakeman-to-codequality:latest > codequality.json
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

### GitLab CI

This repository includes a `.gitlab-ci.yml` that automatically:
- Builds the Docker image and pushes it to GitLab Container Registry
- Runs tests on the image
- Demonstrates usage with Brakeman scanning

The image is available at: `$CI_REGISTRY_IMAGE:latest` (for the main branch) or `$CI_REGISTRY_IMAGE:$CI_COMMIT_REF_SLUG` (for specific branches/tags)

#### Using the built image in your project:

```yaml
brakeman-scan:
  stage: security
  image: ruby:3.2
  services:
    - docker:24-dind
  variables:
    # Replace with your registry URL
    CONVERTER_IMAGE: registry.gitlab.com/your-namespace/brakeman-to-codequality:latest
  before_script:
    - gem install brakeman
  script:
    - brakeman -f json | docker run --rm -i ${CONVERTER_IMAGE} > codequality.json
  artifacts:
    reports:
      codequality: codequality.json
```

#### Alternative: Using Go installation

```yaml
brakeman-scan:
  stage: security
  image: ruby:3.2
  before_script:
    - apt-get update && apt-get install -y golang-go
    - gem install brakeman
    - go install github.com/Omochice/brakeman-to-codequality@latest
    - export PATH=$PATH:$(go env GOPATH)/bin
  script:
    - brakeman -f json | brakeman-to-codequality > codequality.json
  artifacts:
    reports:
      codequality: codequality.json
```

See `.gitlab-ci.yml.example` for more configuration patterns.

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

MIT

## Contributing

Contributions are welcome! Please open an issue or pull request.
