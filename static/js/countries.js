document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.querySelector("#country-search");
    const regionFilter = document.querySelector("#region-filter");
    const resultsContainer = document.querySelector("#country-results");

    if (!searchInput || !regionFilter || !resultsContainer) {
      return;
    }

    let debounceTimer = null;

    function debounce(callback, delay = 350) {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(callback, delay);
    }

    function escapeHTML(value) {
      if (value === null || value === undefined) {
        return "";
      }

      return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
    }

    function formatLanguages(languages) {
      if (!Array.isArray(languages) || languages.length === 0) {
        return "N/A";
      }

      return languages.map(escapeHTML).join(", ");
    }

    function renderLoading() {
      resultsContainer.innerHTML = `
        <div class="empty-panel">
          Loading countries...
        </div>
      `;
    }

    function renderError(message) {
      resultsContainer.innerHTML = `
        <div class="alert-panel">
          ${escapeHTML(message)}
        </div>
      `;
    }

    function renderEmpty() {
      resultsContainer.innerHTML = `
        <div class="empty-panel">
          No countries found.
        </div>
      `;
    }

    function renderCountries(countries) {
      if (!Array.isArray(countries) || countries.length === 0) {
        renderEmpty();
        return;
      }

      resultsContainer.innerHTML = countries
        .map((country) => {
          const name = escapeHTML(country.name);
          const slug = escapeHTML(country.slug);
          const flag = escapeHTML(country.flag);
          const official = escapeHTML(country.official || "");
          const region = escapeHTML(country.region || "N/A");
          const capital = escapeHTML(country.capital || "N/A");
          const population = escapeHTML(country.population || "N/A");
          const currency = escapeHTML(country.currency || "N/A");

          return `
            <a href="/countries/${slug}" class="country-card">
                <div class="country-card-media">
                    ${flag ? `<img src="${flag}" alt="${name} flag">` : ""}
                </div>

                <div class="country-card-body">
                    <h2>${name}</h2>
                    <p>Capital: ${capital}</p>
                    <p>Population: ${population}</p>
                    <p>Currency: ${currency}</p>
                    <p>Languages: ${formatLanguages(country.languages)}</p>
                </div>
            </a>
            `;
        })
        .join("");
    }

    async function loadCountries() {
      const params = new URLSearchParams();

      const search = searchInput.value.trim();
      const region = regionFilter.value.trim();

      if (search !== "") {
        params.append("search", search);
      }

      if (region !== "") {
        params.append("region", region);
      }

      renderLoading();

      try {
        const response = await fetch(`/api/countries?${params.toString()}`, {
          headers: {
            Accept: "application/json",
          },
        });

        const data = await response.json();

        if (!response.ok) {
          renderError(data.message || "Could not load countries.");
          return;
        }

        renderCountries(data);
      } catch (error) {
        renderError("Network error. Please try again.");
      }
    }

    searchInput.addEventListener("input", () => {
      debounce(loadCountries);
    });

    regionFilter.addEventListener("change", loadCountries);
  });