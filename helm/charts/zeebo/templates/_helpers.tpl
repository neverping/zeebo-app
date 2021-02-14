{{/* vim: set filetype=mustache: */}}

{{/*
  Common annotations to be used in all resources
*/}}
{{- define "zeebo.annotations" -}}
team: {{ .Values.team }}
repository: {{ .Values.repository }}
{{- end -}}

{{/*
  Common labels to be used in all resources
  More info:
  - https://helm.sh/docs/chart_best_practices/labels/
  - https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
*/}}
{{- define "zeebo.labelsCommon" -}}
helm.sh/chart: {{ .Chart.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/part-of: {{ .Chart.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
  Selector labels Go
*/}}
{{- define "zeebo.selectorLabelsGo" -}}
app.kubernetes.io/name: {{ .Values.apps.go.svc.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
  Selector labels Python
*/}}
{{- define "zeebo.selectorLabelsPython" -}}
app.kubernetes.io/name: {{ .Values.apps.python.svc.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
  Common labels fo Go App
*/}}
{{- define "zeebo.labelsGo" -}}
{{ include "zeebo.labelsCommon" . }}
{{ include "zeebo.selectorLabelsGo" . }}
{{- end -}}

{{/*
  Common labels fo Python App
*/}}
{{- define "zeebo.labelsPython" -}}
{{ include "zeebo.labelsCommon" . }}
{{ include "zeebo.selectorLabelsPython" . }}
{{- end -}}