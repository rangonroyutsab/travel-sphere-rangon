<section class="page-shell container">
    <div class="page-heading">
        <p class="eyebrow">Personal Travel List</p>
        <h1>My Wishlist</h1>
        <p>Edit notes, update trip status, or remove saved destinations.</p>
    </div>
    <section class="wishlist-panel">
        <div id="wishlist-page-feedback" class="feedback-area"></div>
        <div class="wishlist-table-wrap">
            <table class="wishlist-table">
                <thead>
                    <tr>
                        <th>Country</th>
                        <th>Note</th>
                        <th>Status</th>
                        <th>Created</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody id="wishlist-rows">
                    {{range .WishlistItems}}
                    <tr data-id="{{.ID}}">
                        <td>
                            <strong>{{.CountryName}}</strong>
                        </td>
                        <td>
                            <input class="note-input" type="text" value="{{.Note}}" placeholder="Add a note...">
                        </td>
                        <td>
                            <select class="status-select">
                            <option value="Planned" {{if eq .Status "Planned"}}selected{{end}}>Planned</option>
                            <option value="Visited" {{if eq .Status "Visited"}}selected{{end}}>Visited</option>
                            </select>
                        </td>
                        <td>
                            <span class="muted">{{.CreatedAt}}</span>
                        </td>
                        <td>
                            <div class="row-actions">
                                <button class="save-btn" type="button">Save</button>
                                <button class="delete-btn" type="button">Delete</button>
                            </div>
                        </td>
                    </tr>
                    {{else}}
                    <tr>
                        <td colspan="5" class="empty-state">
                            Your wishlist is empty. Explore countries and add destinations.
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
    </section>
</section>
<script src="/static/js/wishlist.js?v=1"></script>
