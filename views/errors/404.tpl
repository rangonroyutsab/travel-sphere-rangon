  <section class="page">
    <h1>{{.Title}}</h1>

    {{if .Message}}
      <p>{{.Message}}</p>
    {{else}}
      <p>The page you requested could not be found.</p>
    {{end}}

    <a href="/countries">Back to countries</a>
  </section>