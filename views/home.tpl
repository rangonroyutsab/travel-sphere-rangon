<section class="hero-section">
    <div class="subtle-pattern"></div>

    <div class="hero-content container">
        <p class="eyebrow">Global Country Explorer</p>

        <h1>
            Discover your next
            <span>destination</span>
        </h1>

        <p class="hero-copy">
            Search countries, compare regions, and open destination pages with useful country facts.
        </p>

        <div class="hero-search">
            <div class="search-box">
                <span class="search-icon">⌕</span>
                <input id="home-search" type="text" placeholder="Search countries..." autocomplete="off">
                <a href="/countries" class="primary-button">Explore</a>
            </div>

            <div id="search-suggestions" class="search-suggestions">
                <p>Search suggestions will appear here.</p>
            </div>
        </div>
    </div>
</section>

<section class="section container">
    <div class="section-heading">
        <div>
            <p class="eyebrow">Start Exploring</p>
            <h2>Featured destinations</h2>
        </div>

        <a href="/countries" class="text-link">View all</a>
    </div>

    {{if .CountryError}}
        <div class="alert-panel">
            {{.CountryError}}
        </div>
    {{else}}
        <div class="featured-grid">
            {{range .FeaturedCountries}}
            <a href="/countries/{{.Slug}}" class="featured-card">
                {{if .Flag}}
                    <img class="flag-img" src="{{.Flag}}" alt="{{.Name}} flag">
                {{else}}
                    <span class="flag">{{.Region}}</span>
                    {{end}}

                    <span class="badge">{{.Region}}</span>

                    <h3>{{.Name}}</h3>

                    <p>{{.Capital}} · {{.Population}}</p>
                </a>
            {{else}}
                <div class="empty-panel">
                    No featured destinations are available right now.
                </div>
                {{end}}
            </div>
            {{end}}
        </section>