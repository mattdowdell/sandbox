# https://docs.tilt.dev/api.html

# -----
# Setup
# -----

# Prevent tilt from accessing other clusters
allow_k8s_contexts("kind-kind")

# ----------
# Extensions
# ----------

load("ext://helm_remote", "helm_remote")

# --------------
# Local services
# --------------

# TODO

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
    version="0.27.0",
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
    version="0.5.0",
    values=[".tilt/postgresql/values.yaml"],
)

k8s_resource(
    "cluster-postgresql",
    new_name="postgresql",
    labels="postgresql",
    resource_deps=["cloudnative-pg"],
    port_forwards=["127.0.0.1:5432:5432"],
)

# -------------
# Observability
# -------------

helm_remote(
    "tempo",
    repo_name="grafana",
    repo_url="https://grafana.github.io/helm-charts",
    # renovate: datasource=helm depName=tempo packageName=tempo registryUrl=https://grafana.github.io/helm-charts
    version="1.24.1",
)

k8s_resource(
    "tempo",
    objects=[
        "tempo:serviceaccount",
        "tempo:configmap",
    ],
    labels=["observability"],
)

helm_remote(
    "victoria-metrics-single",
    repo_name="vm",
    repo_url="https://victoriametrics.github.io/helm-charts",
    # renovate: datasource=helm depName=victoria-metrics-single packageName=victoria-metrics-single registryUrl=https://victoriametrics.github.io/helm-charts
    version="0.27.0",
)

k8s_resource(
    "victoria-metrics-single-server",
    objects=["victoria-metrics-single-server:serviceaccount"],
    labels=["observability"],
    port_forwards=["127.0.0.1:8428:8428"],
)

helm_remote(
    "victoria-logs-single",
    repo_name="vm",
    repo_url="https://victoriametrics.github.io/helm-charts",
    # renovate: datasource=helm depName=victoria-logs-single packageName=victoria-logs-single registryUrl=https://victoriametrics.github.io/helm-charts
    version="0.11.23",
)

k8s_resource(
    "victoria-logs-single-server",
    labels=["observability"],
    port_forwards=["127.0.0.1:9428:9428"],
)

helm_remote(
    "grafana",
    repo_name="grafana",
    repo_url="https://grafana.github.io/helm-charts",
    # renovate: datasource=helm depName=grafana packageName=grafana registryUrl=https://grafana.github.io/helm-charts
    version="10.5.5",
    values=[".tilt/grafana/values.yaml"],
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
    ],
    labels=["observability"],
    port_forwards=["127.0.0.1:3000:3000"],
    resource_deps=[
        "tempo",
        "victoria-metrics-single-server",
        "victoria-logs-single-server",
    ],
)
