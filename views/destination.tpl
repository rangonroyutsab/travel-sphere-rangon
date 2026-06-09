<section class="page-shell container">
    <a href="/countries" class="text-link">Back to countries</a>
    <section class="destination-hero">
        <div class="destination-flag-panel">
            {{if .Country.Flag}}
            <img src="{{.Country.Flag}}" alt="{{.Country.Name}} flag">
            {{end}}
        </div>
        <div class="destination-summary">
            <p class="eyebrow">{{.Country.Region}}</p>
            <h1>{{.Country.Name}}</h1>
            <p>{{.Country.Official}}</p>
        </div>
    </section>
    <section class="wishlist-action-panel">
        {{if .IsAuthenticated}}
        <button id="add-wishlist-btn" class="primary-button" type="button" data-country="{{.Country.Name}}">
        Add to Wishlist
        </button>
        {{else}}
        <a href="/login" class="primary-button">Login to add this destination</a>
        {{end}}
        <div id="wishlist-feedback" class="feedback-area"></div>
    </section>
    <section class="detail-grid">
        <article class="detail-card">
            <span>Capital</span>
            <strong>{{.Country.Capital}}</strong>
        </article>
        <article class="detail-card">
            <span>Population</span>
            <strong>{{.Country.Population}}</strong>
        </article>
        <article class="detail-card">
            <span>Region</span>
            <strong>{{.Country.Region}}</strong>
        </article>
        <article class="detail-card">
            <span>Subregion</span>
            <strong>{{.Country.Subregion}}</strong>
        </article>
        <article class="detail-card">
            <span>Currency</span>
            <strong>{{.Country.Currency}}</strong>
        </article>
        <article class="detail-card">
            <span>Coordinates</span>
            <strong>{{.Country.Lat}}, {{.Country.Lng}}</strong>
        </article>
    </section>
    <section class="detail-panel">
        <div>
            <p class="eyebrow">Languages</p>
            <p>
                {{range $i, $lang := .Country.Languages}}
                {{if $i}}, {{end}}{{$lang}}
                {{else}}
                N/A
                {{end}}
            </p>
        </div>
    </section>
    <section class="detail-panel">
        <div>
            <p class="eyebrow">Attractions</p>
            <h2>Coming later</h2>
            <p>
                Attractions will be added in a later feature slice using OpenTripMap.
            </p>
        </div>
    </section>
</section>
<script src="/static/js/wishlist.js?v=1"></script>
