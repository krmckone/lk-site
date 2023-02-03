{{range makeHrefs "assets/pages/posts"}}
<a class="navbar-item" href="{{.}}.html">
  {{makeNavTitle .}}
</a>
{{end}}
