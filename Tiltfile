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

# TODO

# -------------
# Observability
# -------------

helm_remote(
    "tempo",
    repo_name="grafana",
    repo_url="https://grafana.github.io/helm-charts",
    # renovate: datasource=helm depName=tempo packageName=tempo registryUrl=grafana.github.io/helm-charts
    version="1.23.3",
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
    # renovate: datasource=helm depName=victoria-metrics-single packageName=victoria-metrics-single registryUrl=victoriametrics.github.io/helm-charts
    version="0.24.4",
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
    # renovate: datasource=helm depName=victoria-logs-single packageName=victoria-logs-single registryUrl=victoriametrics.github.io/helm-charts
    version="0.11.7",
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
    # renovate: datasource=helm depName=grafana packageName=grafana registryUrl=grafana.github.io/helm-charts
    version="9.4.4",
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
