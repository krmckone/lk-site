{{range getAssets "assets/pages"}}
<a class="navbar-item" href="{{.}}.html">
  {{makeTitle .}}
</a>
{{end}}