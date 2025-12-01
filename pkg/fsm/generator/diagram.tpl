@startuml

<style>
componentDiagram {
  BackGroundColor white
  LineThickness 1
  LineColor #333333
}
document {
  BackGroundColor white
}
</style>

skinparam defaulttextalignment center

top to bottom direction
hide empty description
{{- $colourStart := .Start }}
{{- $colourStop := .Stop }}
{{- $colourState := .State }}
{{- $colourLink := .Link }}
{{ range $state := .States }}
state state_{{ $state.Name }} as "<b>{{ $state.Name }}</b>"#{{ if $state.Colour }}{{ trimHash $state.Colour }}{{ else }}{{ trimHash $colourState }}{{ end }}{{ if $state.Description }}:{{ $state.Description }}{{ end }}
{{- end }}

{{ if .Initial }}state {{ .Initial }}_start <<start>>#{{ trimHash $colourStart }}
{{ .Initial }}_start -[#{{ trimHash $colourLink }}]-> state_{{ .Initial }}{{ end }}
{{ range $edge := .Edges }}
state_{{ $edge.From }} -[#{{ if $edge.Colour }}{{ trimHash $edge.Colour }}{{ else }}{{ trimHash $colourLink }}{{ end }}]-> state_{{ $edge.To }}{{ if $edge.Description }}: {{ $edge.Description }}{{ end }}
{{- end }}
{{ range $final := .Final }}
state {{ $final.Name }}_end <<end>>#{{ trimHash $colourStop }}
state_{{ $final.Name }} -[#{{ trimHash $colourLink }}]-> {{ $final.Name }}_end
{{- end }}

@enduml
