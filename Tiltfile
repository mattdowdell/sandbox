# https://docs.tilt.dev/api.html

# -----
# Setup
# -----

# Prevent tilt from accessing other clusters
allow_k8s_contexts("kind-example")

# Some charts seems to time out with the default of 30 seconds.
update_settings(k8s_upsert_timeout_secs=120)

# ----------
# Extensions
# ----------

load("ext://helm_resource", "helm_repo", "helm_resource")

# --------------
# Local services
# --------------

# Tilt's docker_build excludes .git per https://github.com/tilt-dev/tilt/issues/2169
# but we need that for db migrations and otel resource attributes
def custom_docker_build(name, target="runtime", go_build_args=""):
    command = "docker buildx build --pull \
        --target {target} \
        --build-arg SERVICE={name} \
        --build-arg GO_BUILD_ARGS={go_build_args!r} \
        -t $EXPECTED_REF ." \
        .format(
            name=name,
            target=target,
            go_build_args=go_build_args,
        )

    custom_build(
        ref=name,
        command=command,
        deps=[
            os.path.join(config.main_dir, "Dockerfile"),
            os.path.join(config.main_dir, "cmd"),
            os.path.join(config.main_dir, "internal"),
            os.path.join(config.main_dir, "pkg"),
        ],
    )

go_build_args = os.getenv("GO_BUILD_ARGS", "")
custom_docker_build("example-rpc", target="debug", go_build_args=go_build_args)
custom_docker_build("example-db-migrate", go_build_args=go_build_args)

k8s_yaml(kustomize(os.path.join(config.main_dir, "k8s/kustomize/example/tilt")))

k8s_resource(
    "example-rpc",
    objects=[
        "example:configmap",
        "example:serviceaccount",
        "example-rpc:poddisruptionbudget",
    ],
    labels=["services"],
    port_forwards=["127.0.0.1:5000:5000"],
    resource_deps=[
        "postgresql-cluster",
        "victoria-traces",
        "victoria-metrics",
        "victoria-logs",
    ],
)

k8s_resource(
    "example-db-migrate",
    labels=["services"],
    resource_deps=[
        "postgresql-cluster",
        "victoria-traces",
        "victoria-metrics",
        "victoria-logs",
    ],
)

# --------
# Database
# --------

helm_repo(
    name="cnpg",
    url="https://cloudnative-pg.github.io/charts",
    resource_name="cnpg-repo",
    labels=["postgresql"],
)

helm_resource(
    name="cloudnative-pg",
    chart="cnpg/cloudnative-pg",
    namespace="cnpg-system",
    flags=[
        "--create-namespace",
        "--version=0.28.0",
    ],
    resource_deps=["cnpg-repo"],
    labels=["postgresql"],
)

k8s_kind(
    "Cluster",
    api_version="postgresql.cnpg.io/v1",
    image_json_path="{.spec.imageName}",
)

helm_resource(
    name="postgresql-cluster",
    chart="cnpg/cluster",
    flags=[
        # renovate: datasource=helm depName=cluster packageName=cluster registryUrl=https://cloudnative-pg.github.io/charts
        "--version=0.6.1",
        "--values=k8s/helm/postgresql/values.yaml",
    ],
    deps=["k8s/helm/postgresql/values.yaml"],
    resource_deps=[
        "cnpg-repo",
        "cloudnative-pg",
    ],
    labels=["postgresql"],
    port_forwards=["127.0.0.1:5432:5432"],
)

# ------------
# Cert Manager
# ------------

helm_repo(
    "jetstack",
    "https://charts.jetstack.io", # TODO: migrate to oci repo
    resource_name="cert-manager-repo",
    labels=["cert-manager"],
)

helm_resource(
    "cert-manager",
    "jetstack/cert-manager",
    resource_deps=["cert-manager-repo"],
    labels=["cert-manager"],
    namespace="cert-manager",
    flags=[
        # renovate: datasource=helm depName=cert-manager packageName=cert-manager registryUrl=https://charts.jetstack.io
        "--version=1.20.2",
        "--create-namespace",
        "--values=k8s/helm/cert-manager/values.yaml",
    ],
)

# -------------
# OpenTelemetry
# -------------

helm_repo(
    "open-telemetry",
    "https://open-telemetry.github.io/opentelemetry-helm-charts",
    resource_name="opentelemetry-repo",
    labels="opentelemetry",
)

helm_resource(
    "opentelemetry-operator",
    "open-telemetry/opentelemetry-operator",
    namespace="opentelemetry",
    labels=["opentelemetry"],
    flags=[
        "--create-namespace",
        # renovate: datasource=helm depName=opentelemetry-operator packageName=opentelemetry-operator registryUrl=https://open-telemetry.github.io/opentelemetry-helm-charts
        "--version=0.115.0",
        "--values=k8s/helm/opentelemetry-operator/values.yaml",
    ],
    resource_deps=["cert-manager"],
)

# There's a delay between the operator deployment becoming ready
# and the mutating webhook for OpentelemetryCollector CRD being usable
# https://github.com/open-telemetry/opentelemetry-operator/issues/3194
local_resource(
    "opentelemetry-operator-ready",
    "./tools/wait-for-otel-operator.sh",
    resource_deps=["opentelemetry-operator"],
    labels=["opentelemetry"],
)

k8s_kind(
    "OpenTelemetryCollector",
    api_version="opentelemetry.io/v1beta1",
)

k8s_yaml(kustomize(os.path.join(config.main_dir, "k8s/kustomize/opentelemetry-collector/base")))

k8s_resource(
    "opentelemetry-daemonset",
    labels=["opentelemetry"],
    resource_deps=[
        "opentelemetry-operator-ready",
    ],
)

# -------------
# Observability
# -------------

helm_repo(
    name="vm",
    url="https://victoriametrics.github.io/helm-charts",
    resource_name="vm-repo",
    labels=["observability"],
)

helm_repo(
    name="grafana-community",
    url="https://grafana-community.github.io/helm-charts",
    resource_name="grafana-repo",
    labels=["observability"],
)

helm_resource(
    name="victoria-traces",
    chart="vm/victoria-traces-single",
    namespace="observability",
    flags=[
        "--create-namespace",
        # renovate: datasource=helm depName=victoria-traces-single packageName=victoria-traces-single registryUrl=https://victoriametrics.github.io/helm-charts
        "--version=0.1.8",
    ],
    resource_deps=["vm-repo"],
    labels=["observability"],
    port_forwards=["127.0.0.1:10428:10428"],
)

helm_resource(
    name="victoria-metrics",
    chart="vm/victoria-metrics-single",
    namespace="observability",
    flags=[
        "--create-namespace",
        # renovate: datasource=helm depName=victoria-metrics-single packageName=victoria-metrics-single registryUrl=https://victoriametrics.github.io/helm-charts
        "--version=0.39.0",
    ],
    resource_deps=["vm-repo"],
    labels=["observability"],
    port_forwards=["127.0.0.1:8428:8428"],
)

helm_resource(
    name="victoria-logs",
    chart="vm/victoria-logs-single",
    namespace="observability",
    flags=[
        "--create-namespace",
        # renovate: datasource=helm depName=victoria-logs-single packageName=victoria-logs-single registryUrl=https://victoriametrics.github.io/helm-charts
        "--version=0.13.7",
    ],
    resource_deps=["vm-repo"],
    labels=["observability"],
    port_forwards=["127.0.0.1:9428:9428"],
)

k8s_yaml(kustomize(os.path.join(config.main_dir, "k8s/kustomize/grafana/base")))

k8s_resource(
    new_name="grafana-dashboards",
    objects=["grafana-dashboards:configmap"],
    labels=["observability"],
)

helm_resource(
    name="grafana",
    chart="grafana-community/grafana",
    namespace="observability",
    flags=[
        "--create-namespace",
        # renovate: datasource=helm depName=grafana packageName=grafana registryUrl=https://grafana.github.io/helm-charts
        "--version=10.5.15",
        "--values=k8s/helm/grafana/values.yaml"
    ],
    deps=["k8s/helm/grafana/values.yaml"],
    resource_deps=["grafana-dashboards", "vm-repo"],
    labels=["observability"],
    port_forwards=["127.0.0.1:3000:3000"],
)
