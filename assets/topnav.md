{{range getAssetsNoAbout "assets/pages"}}
<a class="navbar-item" href="{{.}}.html">
  {{makeTitle .}}
</a>
{{end}}