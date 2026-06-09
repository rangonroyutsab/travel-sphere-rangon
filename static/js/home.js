document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.querySelector("#home-search");
    const suggestionsContainer = document.querySelector("#search-suggestions");

    if (!searchInput || !suggestionsContainer) {
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

    function renderDefault() {
        suggestionsContainer.innerHTML = `
        <p>Search suggestions will appear here.</p>
      `;
    }

    function renderLoading() {
        suggestionsContainer.innerHTML = `
        <p>Searching countries...</p>
      `;
    }

    function renderError(message) {
        suggestionsContainer.innerHTML = `
        <div class="alert-panel">
          ${escapeHTML(message)}
        </div>
      `;
    }

    function renderSuggestions(countries) {
        if (!Array.isArray(countries) || countries.length === 0) {
            suggestionsContainer.innerHTML = `
          <p>No countries found.</p>
        `;
            return;
        }

        suggestionsContainer.innerHTML = countries
            .slice(0, 6)
            .map((country) => {
                const name = escapeHTML(country.name);
                const slug = escapeHTML(country.slug);
                const capital = escapeHTML(country.capital || "N/A");
                const region = escapeHTML(country.region || "N/A");

                return `
            <a href="/countries/${slug}" class="suggestion-item">
              <strong>${name}</strong>
              <span>${capital} · ${region}</span>
            </a>
          `;
            })
            .join("");
    }

    async function loadSuggestions() {
        const search = searchInput.value.trim();

        if (search.length < 2) {
            renderDefault();
            return;
        }

        const params = new URLSearchParams();
        params.append("search", search);

        renderLoading();

        try {
            const response = await fetch(`/api/countries?${params.toString()}`, {
                headers: {
                    Accept: "application/json",
                },
            });

            const data = await response.json();

            if (!response.ok) {
                renderError(data.message || "Could not load suggestions.");
                return;
            }

            renderSuggestions(data);
        } catch (error) {
            renderError("Network error. Please try again.");
        }
    }

    searchInput.addEventListener("input", () => {
        debounce(loadSuggestions);
    });
});