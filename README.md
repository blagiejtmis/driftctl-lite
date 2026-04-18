# driftctl-lite

> Lightweight CLI to detect config drift between live cloud resources and IaC definitions

---

## Installation

```bash
go install github.com/yourorg/driftctl-lite@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourorg/driftctl-lite/releases).

---

## Usage

Run a drift check against your Terraform state and live AWS resources:

```bash
driftctl-lite scan --provider aws --state ./terraform.tfstate
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--provider` | Cloud provider to scan (`aws`, `gcp`, `azure`) |
| `--state` | Path to your IaC state file |
| `--output` | Output format: `text` (default), `json`, `yaml` |
| `--region` | Cloud region to target |

### Example Output

```
[DRIFT DETECTED]
  Resource: aws_s3_bucket.my-bucket
  Field:    versioning.enabled
  Expected: true
  Actual:   false

Summary: 3 resources scanned, 1 drift(s) found.
```

---

## Requirements

- Go 1.21+
- Cloud provider credentials configured (e.g., AWS credentials via `~/.aws/credentials` or environment variables)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

This project is licensed under the [MIT License](LICENSE).