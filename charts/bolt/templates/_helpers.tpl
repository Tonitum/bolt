{{- define "bolt.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "bolt.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "bolt.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "bolt.labels" -}}
app.kubernetes.io/name: {{ include "bolt.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
