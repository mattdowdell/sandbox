{{/*
Expand the name of the chart.
*/}}
{{- define "helpers.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "helpers.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "helpers.fullname" -}}
{{- if .name -}}
{{- if .context.Values.releasePrefix -}}
{{- printf "%s-%s" .context.Values.releasePrefix .name | trunc 63 | trimAll "-" -}}
{{- else -}}
{{- printf "%s-%s" (include "helpers.name" .context) .name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- else -}}
{{- include "helpers.name" .context -}}
{{- end -}}
{{- end -}}

{{- define "helpers.labels" -}}
helm.sh/chart: {{ include "helpers.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Chart.AppVersion }}
version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- end }}

{{- define "helpers.selectorLabels" -}}
{{- $ := .context -}}
app: {{ include "helpers.fullname" (dict "name" .name "context" $) }}
app.kubernetes.io/name: {{ include "helpers.fullname" (dict "name" .name "context" $) }}
app.kubernetes.io/instance: {{ include "helpers.name" $ }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "helpers.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "helpers.name" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
