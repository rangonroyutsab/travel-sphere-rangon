<section class="page-shell">
  <div class="container dashboard-layout">

    <section class="dashboard-stats">
      <article class="stat-card">
        <span>Total Saved</span>
        <strong>{{.TotalSaved}}</strong>
      </article>

      <article class="stat-card">
        <span>Planned</span>
        <strong>{{.Planned}}</strong>
      </article>

      <article class="stat-card">
        <span>Visited</span>
        <strong>{{.Visited}}</strong>
      </article>
    </section>

    <section class="dashboard-trips">
      <p class="eyebrow">Saved Destinations</p>
      <h1>Your trips</h1>

      {{if .SavedDestinations}}
        <div class="trip-list">
          {{range .SavedDestinations}}
            <article class="trip-card">
              <h3>{{.CountryName}}</h3>
              <p>{{.Status}}</p>
            </article>
          {{end}}
        </div>
      {{else}}
        <p class="muted">You have not saved any destinations yet.</p>
      {{end}}
    </section>

  </div>
</section>