#!/usr/bin/env bash
# Wait for the cloudnative-pg Operator to be available and verify it can create a Cluster CR

set -u

TIMEOUT=300
POLL_INTERVAL=0.5

usage() {
  cat <<EOF
Usage: $0 [--timeout seconds]

Options:
  --timeout Timeout in seconds for each check (default: $TIMEOUT)
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --timeout)
      TIMEOUT="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown arg: $1" >&2
      usage
      exit 2
      ;;
  esac
done

end_time=$(( $(date +%s) + TIMEOUT ))

echo "Waiting until the cloudnative-pg Operator deployment is created"
until [[ $(date +%s) -ge $end_time ]]; do
  if kubectl -n cnpg-system get deployment cloudnative-pg >/dev/null 2>&1; then
    break
  fi
  sleep "$POLL_INTERVAL"
done

if ! kubectl -n cnpg-system get deployment cloudnative-pg >/dev/null 2>&1; then
  echo "Timed out waiting for cloudnative-pg Operator deployment (proceeding to CR check)" >&2
else
  echo "cloudnative-pg Operator deployment is created. Now checking if it is fully operational."
fi

# Prepare a minimal OpenTelemetryCollector CR manifest
read -r -d '' COLLECTOR_YAML <<'EOF' || true
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: operator-check
spec:
  instances: 1
  storage:
    size: 1Gi
EOF

echo "Check if the Cluster CR can be created."
end_time=$(( $(date +%s) + TIMEOUT ))
created=false
while [[ $(date +%s) -lt $end_time ]]; do
  # Try to create the CR; if it fails, sleep and retry. We ignore stderr here because failures are expected
  echo "$COLLECTOR_YAML" | kubectl create -f - >/dev/null 2>&1 && created=true && break
  sleep "$POLL_INTERVAL"
done

if ! $created; then
  echo "Timed out creating Cluster CR" >&2
  exit 1
fi

# Clean up the test CR
if ! kubectl delete cluster operator-check -n default >/dev/null 2>&1; then
  echo "Failed to delete Cluster CR" >&2
  exit 1
fi

echo "cloudnative-pg operator is ready."
