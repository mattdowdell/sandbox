#!/usr/bin/env bash
# Wait for the OpenTelemetry Operator to be available and verify it can create a collector CR
# 
# This is a bash port of https://github.com/open-telemetry/opentelemetry-operator/blob/main/hack/check-operator-ready.go
#
# This uses bash instead of Go to avoid various k8s dependencies which can be painful to manage.
# Specifically, when 1 dependency wants to upgrade k8s, but another does not. Or if a dependency
# gets into the weeds with k8s's module replace directives

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

echo "Waiting until the OpenTelemetry Operator deployment is created"
until [[ $(date +%s) -ge $end_time ]]; do
  if kubectl -n opentelemetry get deployment opentelemetry-operator >/dev/null 2>&1; then
    break
  fi
  sleep "$POLL_INTERVAL"
done

if ! kubectl -n opentelemetry get deployment opentelemetry-operator >/dev/null 2>&1; then
  echo "Timed out waiting for OpenTelemetry Operator deployment (proceeding to CR check)" >&2
else
  echo "OpenTelemetry Operator deployment is created. Now checking if it is fully operational."
fi

# Prepare a minimal OpenTelemetryCollector CR manifest
read -r -d '' COLLECTOR_YAML <<'EOF' || true
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: operator-check
  namespace: default
spec: {}
EOF

echo "Check if the OpenTelemetry collector CR can be created."
end_time=$(( $(date +%s) + TIMEOUT ))
created=false
while [[ $(date +%s) -lt $end_time ]]; do
  # Try to create the CR; if it fails, sleep and retry. We ignore stderr here because failures are expected
  echo "$COLLECTOR_YAML" | kubectl create -f - >/dev/null 2>&1 && created=true && break
  sleep "$POLL_INTERVAL"
done

if ! $created; then
  echo "Timed out creating OpenTelemetry collector CR" >&2
  exit 1
fi

# Clean up the test CR
if ! kubectl delete OpenTelemetryCollector operator-check -n default >/dev/null 2>&1; then
  echo "Failed to delete OpenTelemetry collector CR" >&2
  exit 1
fi

echo "OpenTelemetry operator is ready."
