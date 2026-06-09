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
                <div class="featured-card-media">
                    {{if .Flag}}
                        <img src="{{.Flag}}" alt="{{.Name}} flag">
                        {{end}}
                    </div>

                    <div class="featured-card-body">
                        <h3>{{.Name}}</h3>
                        <p>{{.Capital}} · {{.Region}}</p>
                    </div>
                </a>
            {{else}}
                <div class="empty-panel">
                    No featured destinations are available right now.
                </div>
                {{end}}
            </div>
            {{end}}
        </section>

        <script src="/static/js/home.js?v=1"></script>