{{- define "app-version" -}}
{{- if $.Values.global.appVersionOverride }}
{{- $.Values.global.appVersionOverride }}
{{- else }}
{{- $.Chart.AppVersion }}
{{- end }}
{{- end }}

{{- define "image.builder" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}builder:{{ include "app-version" $ }}
{{- end }}

{{- define "image.controller" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}controller:{{ include "app-version" $ }}
{{- end }}

{{- define "image.dashboard" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}dashboard:{{ include "app-version" $ }}
{{- end }}

{{- define "image.gateway" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}gateway:{{ include "app-version" $ }}
{{- end }}

{{- define "image.gitea-integration" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}gitea-integration:{{ include "app-version" $ }}
{{- end }}

{{- define "image.sablier" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}sablier:{{ include "app-version" $ }}
{{- end }}

{{- define "image.ssgen" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}ssgen:{{ include "app-version" $ }}
{{- end }}

{{- define "image.migrate" -}}
{{- $.Values.global.image.repository }}{{ $.Values.global.image.prefix }}migrate:{{ include "app-version" $ }}
{{- end }}

{{- define "config-hash" }}
{{- include (print $.Template.BasePath "/config.yaml") $ | sha256sum }}
{{- end }}

{{- define "buildkit-config-hash" }}
{{- include (print $.Template.BasePath "/builder/config.yaml") $ | sha256sum }}
{{- end }}

{{- define "sablier-config-hash" }}
{{- include (print $.Template.BasePath "/sablier/config.yaml") $ | sha256sum }}
{{- end }}

{{- define "known-hosts-hash" }}
{{- include (print $.Template.BasePath "/known-hosts.yaml") $ | sha256sum }}
{{- end }}
