<header class="site-header">
    <nav class="site-header-inner container">
        <a href="/" class="site-logo">TravelSphere</a>

        <div class="site-nav">
            <a href="/">Home</a>
            <a href="/countries">Countries</a>
            <a href="/wishlist">Wishlist</a>
            <a href="/dashboard">Dashboard</a>
        </div>

        <div class="site-actions">
            {{if .IsAuthenticated}}
                <span class="user-chip">{{.UserName}}</span>
                <a href="/logout" class="nav-button">Logout</a>
            {{else}}
                <a href="/login" class="nav-button">Login</a>
                {{end}}
            </div>
        </nav>
    </header>