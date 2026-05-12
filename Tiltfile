# https://docs.tilt.dev/api.html

# -----
# Setup
# -----

# Prevent tilt from accessing other clusters
allow_k8s_contexts("kind-example")

# ----------
# Extensions
# ----------

load("ext://helm_remote", "helm_remote")

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
        "postgresql",
        "victoria-traces",
        "victoria-metrics",
        "victoria-logs",
    ],
)

k8s_resource(
    "example-db-migrate",
    labels=["services"],
    resource_deps=[
        "postgresql",
        "victoria-traces",
        "victoria-metrics",
        "victoria-logs",
    ],
)

# --------
# Database
# --------

helm_remote(
    "cloudnative-pg",
    repo_name="cnpg",
    repo_url="https://cloudnative-pg.github.io/charts",
    namespace="cnpg-system",
    create_namespace=True,
    # renovate: datasource=helm depName=cloudnative-pg packageName=cloudnative-pg registryUrl=https://cloudnative-pg.github.io/charts
    version="0.28.0",
)

k8s_resource(
    "cloudnative-pg",
    objects=[
        "cnpg-system:namespace",
        "backups.postgresql.cnpg.io:customresourcedefinition",
        "clusterimagecatalogs.postgresql.cnpg.io:customresourcedefinition",
        "clusters.postgresql.cnpg.io:customresourcedefinition",
        "databases.postgresql.cnpg.io:customresourcedefinition",
        "failoverquorums.postgresql.cnpg.io:customresourcedefinition",
        "imagecatalogs.postgresql.cnpg.io:customresourcedefinition",
        "poolers.postgresql.cnpg.io:customresourcedefinition",
        "publications.postgresql.cnpg.io:customresourcedefinition",
        "scheduledbackups.postgresql.cnpg.io:customresourcedefinition",
        "subscriptions.postgresql.cnpg.io:customresourcedefinition",
        "cnpg-mutating-webhook-configuration:mutatingwebhookconfiguration",
        "cloudnative-pg:serviceaccount",
        "cloudnative-pg:clusterrole",
        "cloudnative-pg-view:clusterrole",
        "cloudnative-pg-edit:clusterrole",
        "cloudnative-pg:clusterrolebinding",
        "cnpg-controller-manager-config:configmap",
        "cnpg-default-monitoring:configmap",
        "cnpg-validating-webhook-configuration:validatingwebhookconfiguration",
    ],
    labels="postgresql",
)

k8s_kind("Cluster", api_version="postgresql.cnpg.io/v1")

helm_remote(
    "cluster",
    repo_name="cnpg",
    repo_url="https://cloudnative-pg.github.io/charts",
    namespace="default",
    # renovate: datasource=helm depName=cluster packageName=cluster registryUrl=https://cloudnative-pg.github.io/charts
    version="0.6.0",
    values=["k8s/helm/postgresql/values.yaml"],
)

k8s_resource(
    "cluster-postgresql",
    new_name="postgresql",
    objects=[
        "cluster-postgresql-monitoring-logical-replication:configmap"
    ],
    labels="postgresql",
    resource_deps=["cloudnative-pg"],
    port_forwards=["127.0.0.1:5432:5432"],
)

# -------------
# Observability
# -------------

helm_remote(
    "victoria-traces-single",
    repo_name="vm",
    repo_url="https://victoriametrics.github.io/helm-charts",
    # renovate: datasource=helm depName=victoria-traces-single packageName=victoria-traces-single registryUrl=https://victoriametrics.github.io/helm-charts
    version="0.0.7",
)

k8s_resource(
    "victoria-traces-single-vt-single-server",
    new_name="victoria-traces",
    labels=["observability"],
    port_forwards=["127.0.0.1:10428:10428"],
)

helm_remote(
    "victoria-metrics-single",
    repo_name="vm",
    repo_url="https://victoriametrics.github.io/helm-charts",
    # renovate: datasource=helm depName=victoria-metrics-single packageName=victoria-metrics-single registryUrl=https://victoriametrics.github.io/helm-charts
    version="0.38.0",
)

k8s_resource(
    "victoria-metrics-single-server",
    new_name="victoria-metrics",
    objects=["victoria-metrics-single-server:serviceaccount"],
    labels=["observability"],
    port_forwards=["127.0.0.1:8428:8428"],
)

helm_remote(
    "victoria-logs-single",
    repo_name="vm",
    repo_url="https://victoriametrics.github.io/helm-charts",
    # renovate: datasource=helm depName=victoria-logs-single packageName=victoria-logs-single registryUrl=https://victoriametrics.github.io/helm-charts
    version="0.12.4",
)

k8s_resource(
    "victoria-logs-single-server",
    new_name="victoria-logs",
    labels=["observability"],
    port_forwards=["127.0.0.1:9428:9428"],
)

k8s_yaml(kustomize(os.path.join(config.main_dir, "k8s/kustomize/grafana/base")))

helm_remote(
    "grafana",
    repo_name="grafana-community",
    repo_url="https://grafana-community.github.io/helm-charts",
    # renovate: datasource=helm depName=grafana packageName=grafana registryUrl=https://grafana.github.io/helm-charts
    version="10.5.15",
    values=["k8s/helm/grafana/values.yaml"],
)

k8s_resource(
    "grafana",
    objects=[
        "grafana:serviceaccount",
        "grafana:role",
        "grafana-clusterrole:clusterrole",
        "grafana:rolebinding",
        "grafana-clusterrolebinding:clusterrolebinding",
        "grafana:configmap",
        "grafana:secret",
        "grafana-dashboards:configmap",
    ],
    labels=["observability"],
    port_forwards=["127.0.0.1:3000:3000"],
    resource_deps=[
        "victoria-traces",
        "victoria-metrics",
        "victoria-logs",
    ],
)
