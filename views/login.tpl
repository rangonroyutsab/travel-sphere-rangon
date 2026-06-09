<section class="page-shell container">
    <div class="page-heading">
        <p class="eyebrow">User Authentication</p>
        <h1>Login</h1>
        <p>
            Kindly use the provided credentials to log into this demo authentication system.
        </p>
    </div>

    {{if .Error}}
        <p class="form-error">{{.Error}}</p>
    {{end}}

    <div class="login-wrapper">
        <form class="login-form" method="post" action="/login">
            <label class="form-field">
                <span>Username</span>
                <input
                    name="username"
                    type="text"
                    placeholder="Enter username"
                >
            </label>

            <label class="form-field">
                <span>Password</span>
                <input
                    name="password"
                    type="password"
                    placeholder="Enter password"
                >
            </label>

            <button type="submit">Login</button>
        </form>
    </div>
</section>