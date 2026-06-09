<section class="page-shell container">
    <div class="page-heading">
        <p class="eyebrow">Country Explorer</p>
        <h1>Explore destinations</h1>
        <p>
            Browse countries by name, capital, or region. The first results are server-rendered.
        </p>
    </div>

    <section class="filter-panel">
        <label>
            Search
            <input id="country-search" type="text" placeholder="Country, capital, or official name..."
                autocomplete="off">
        </label>

        <label>
            Region
            <select id="region-filter">
                <option value="">All regions</option>
                <option value="Africa">Africa</option>
                <option value="Americas">Americas</option>
                <option value="Asia">Asia</option>
                <option value="Europe">Europe</option>
                <option value="Oceania">Oceania</option>
            </select>
        </label>
    </section>

    {{if .Error}}
        <div class="alert-panel">
            {{.Error}}
        </div>
        {{end}}

        <section id="country-results" class="country-grid">
            {{range .Countries}}
            <a href="/countries/{{.Slug}}" class="country-card">
                <div class="country-card-media">
                    {{if .Flag}}
                        <img src="{{.Flag}}" alt="{{.Name}} flag">
                        {{end}}
                    </div>

                    <div class="country-card-body">
                        <h2>{{.Name}}</h2>

                        <p>Capital: {{.Capital}}</p>
                        <p>Population: {{.Population}}</p>
                        <p>Currency: {{.Currency}}</p>

                        <p>
                            Languages:
                            {{range $i, $lang := .Languages}}
                            {{if $i}}, {{end}}{{$lang}}
                            {{else}}
                                N/A
                                {{end}}
                            </p>
                        </div>
                    </a>
                {{else}}
                    <div class="empty-panel">
                        No countries found.
                    </div>
                    {{end}}
                </section>
            </section>

            <script src="/static/js/countries.js?v=1"></script>