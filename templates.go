package bete

const (
	arrivalsByServiceTemplateString = `{{ if .Stop.Description }}<strong>{{ .Stop.Description }} ({{ .Stop.ID }})</strong>
{{ else }}<strong>{{ .Stop.ID }}</strong>
{{ end }}{{ with .Stop.RoadName }}{{ . }}
{{ end }}<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
{{- range (.Services | filterByService .Filter | sortByService) }}
{{ $fst := until $.Time .NextBus.EstimatedArrival -}}
{{ $snd := until $.Time .NextBus2.EstimatedArrival -}}
{{ $thd := until $.Time .NextBus3.EstimatedArrival -}}
| {{ printf "%-4v" .ServiceNo }} | {{ printf "%3v" $fst }} | {{ printf "%3v" $snd }} | {{ printf "%3v" $thd }} |
{{- end }}
</pre>
{{ with .Filter }}Filtered by services: {{ join . ", " }}
{{ end }}<em>Last updated on {{ .Time | inSGT }}</em>`
)