document.addEventListener("DOMContentLoaded", () => {
    const statsContainer = document.querySelector("#dashboard-stats");

    if (!statsContainer) {
        return;
    }

    refreshDashboardStats(statsContainer);
});

async function refreshDashboardStats(statsContainer) {
    try {
        const response = await fetch("/api/dashboard/summary", {
            headers: {
                Accept: "application/json",
            },
        });

        const data = await response.json();

        if (!response.ok) {
            return;
        }

        renderDashboardStats(statsContainer, data);
    } catch (error) {
        // Keep server-rendered stats if refresh fails.
    }
}

function renderDashboardStats(container, summary) {
    const totalSaved = Number(summary.total_saved || 0);
    const planned = Number(summary.planned || 0);
    const visited = Number(summary.visited || 0);

    container.innerHTML = `
      <article class="stat-card">
        <span>Total Saved</span>
        <strong>${totalSaved}</strong>
      </article>

      <article class="stat-card">
        <span>Planned</span>
        <strong>${planned}</strong>
      </article>

      <article class="stat-card">
        <span>Visited</span>
        <strong>${visited}</strong>
      </article>
    `;
}