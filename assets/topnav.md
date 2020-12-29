{{range getAssets "assets/pages"}}
<a class="navbar-item" href="{{.}}.html">
  {{.}}
</a>
{{end}}