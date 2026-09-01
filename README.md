# write-only-secrets-portal
WOSP: Write-Only Secrets Portal

### Requirements

- mise https://mise.jdx.dev/

### Try or test with moto (simulated AWS SecretsManager API)

1. In a terminal, run `./moto.sh`
2. In another terminal, run `./run-dev.sh`
3. In a browser visit http://localhost:8888/aws/test.html

### Run for real

1. In a terminal, set your AWS credentials environment variables, then run `./run-production.sh`
3. In a browser visit http://localhost:8888/aws/
