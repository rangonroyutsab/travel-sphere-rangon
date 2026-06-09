document.addEventListener("DOMContentLoaded", () => {
    setupAddToWishlist();
    setupWishlistTable();
});

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

function showMessage(container, message, type = "success") {
    if (!container) {
        return;
    }

    const className = type === "error" ? "error-message" : "success-message";

    container.innerHTML = `
      <div class="${className}">
        ${escapeHTML(message)}
      </div>
    `;
}

async function parseJSONResponse(response) {
    const text = await response.text();

    if (!text) {
        return null;
    }

    try {
        return JSON.parse(text);
    } catch (error) {
        return null;
    }
}

function setupAddToWishlist() {
    const addButton = document.querySelector("#add-wishlist-btn");
    const feedback = document.querySelector("#wishlist-feedback");

    if (!addButton || !feedback) {
        return;
    }

    addButton.addEventListener("click", async () => {
        const countryName = addButton.dataset.country;

        if (!countryName) {
            showMessage(feedback, "Country name is missing.", "error");
            return;
        }

        addButton.disabled = true;
        showMessage(feedback, "Adding destination...");

        try {
            const response = await fetch("/api/wishlist", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                },
                body: JSON.stringify({
                    country_name: countryName,
                    note: "",
                    status: "Planned",
                }),
            });

            const data = await parseJSONResponse(response);

            if (!response.ok) {
                showMessage(feedback, data && data.message ? data.message : "Could not add destination.", "error");
                addButton.disabled = false;
                return;
            }

            showMessage(feedback, "Added to wishlist.");
            addButton.textContent = "Added to Wishlist";
        } catch (error) {
            showMessage(feedback, "Network error. Please try again.", "error");
            addButton.disabled = false;
        }
    });
}

function setupWishlistTable() {
    const rowsContainer = document.querySelector("#wishlist-rows");

    if (!rowsContainer) {
        return;
    }

    rowsContainer.addEventListener("click", async (event) => {
        const saveButton = event.target.closest(".save-btn");
        const deleteButton = event.target.closest(".delete-btn");

        if (saveButton) {
            await updateWishlistRow(saveButton.closest("tr"));
            return;
        }

        if (deleteButton) {
            await deleteWishlistRow(deleteButton.closest("tr"));
        }
    });
}

async function updateWishlistRow(row) {
    const feedback = document.querySelector("#wishlist-page-feedback");

    if (!row) {
        return;
    }

    const id = row.dataset.id;
    const noteInput = row.querySelector(".note-input");
    const statusSelect = row.querySelector(".status-select");
    const saveButton = row.querySelector(".save-btn");

    if (!id || !noteInput || !statusSelect || !saveButton) {
        showMessage(feedback, "Wishlist row data is incomplete.", "error");
        return;
    }

    saveButton.disabled = true;

    try {
        const response = await fetch(`/api/wishlist/${encodeURIComponent(id)}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
            body: JSON.stringify({
                note: noteInput.value.trim(),
                status: statusSelect.value,
            }),
        });

        const data = await parseJSONResponse(response);

        if (!response.ok) {
            showMessage(feedback, data && data.message ? data.message : "Could not update wishlist item.", "error");
            saveButton.disabled = false;
            return;
        }

        showMessage(feedback, "Wishlist item updated.");
        saveButton.disabled = false;
    } catch (error) {
        showMessage(feedback, "Network error. Please try again.", "error");
        saveButton.disabled = false;
    }
}

async function deleteWishlistRow(row) {
    const feedback = document.querySelector("#wishlist-page-feedback");

    if (!row) {
        return;
    }

    const id = row.dataset.id;

    if (!id) {
        showMessage(feedback, "Wishlist item id is missing.", "error");
        return;
    }

    try {
        const response = await fetch(`/api/wishlist/${encodeURIComponent(id)}`, {
            method: "DELETE",
            headers: {
                Accept: "application/json",
            },
        });

        const data = await parseJSONResponse(response);

        if (!response.ok) {
            showMessage(feedback, data && data.message ? data.message : "Could not delete wishlist item.", "error");
            return;
        }

        row.remove();
        showMessage(feedback, "Wishlist item deleted.");
        ensureWishlistEmptyState();
    } catch (error) {
        showMessage(feedback, "Network error. Please try again.", "error");
    }
}

function ensureWishlistEmptyState() {
    const rowsContainer = document.querySelector("#wishlist-rows");

    if (!rowsContainer) {
        return;
    }

    if (rowsContainer.querySelectorAll("tr[data-id]").length > 0) {
        return;
    }

    rowsContainer.innerHTML = `
      <tr>
        <td colspan="5" class="empty-state">
          Your wishlist is empty. Explore countries and add destinations.
        </td>
      </tr>
    `;
}