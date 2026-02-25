# AGENTS.md
## Project Structure
```
booruview/
├── api/                    # Go API server (backend)
│   ├── routes/             # HTTP route handlers
│   ├── models/             # SQLC-generated database code
│   └── cmd/main.go       # API server entry point
├── client/               # Vue + TypeScript SPA (frontend)
│   ├── src/              # Vue components & logic
│   │   ├── views/          # Page components (LandingView, SearchResultsView, FavoritesView)
│   │   ├── store.ts        # Centralized state management
│   │   ├── router.ts       # Vue Router configuration
│   │   └── search.ts       # Search query utilities
│   └── vite.config.ts    # Vue build configuration
├── docker-compose.yml      # Dev environment orchestration
└── caddy/                # Caddy web server configuration
```
## Tech Stack
- **Frontend**:
  - Vue 3 + TypeScript
  - Vite build tool
  - Vue Router for navigation
  - Centralized state management via reactive store
## Frontend Code Conventions
- **Vue Single File Components (SFC)**:
  - Components are organized in `src/views/` with dedicated files for each page (e.g., `LandingView.vue`, `SearchResultsView.vue`, `FavoritesView.vue`)
  - SFCs use `<script setup>` syntax for composition API and `<template>` for markup
  - Components interact with the store to access state and dispatch actions
- **State Management**:
  - Centralized `store.ts` using Vue's `reactive()` and `computed()` for reactivity
  - Store manages:
    - User authentication and account data (`account: { authToken, username, data }`)
    - Search state (`query`, `lastQuery`, `posts`, `cachedTags`)
    - UI state (`settings`, `toast`, `sidebarClosed`)
    - Persistent storage via `localStorage` for settings and account data
- **Routing**:
  - Vue Router with `createWebHistory()` for dynamic routing
  - Routes:
    - `/` → `LandingView.vue` (root page)
    - `/search/:page/:query` → `SearchResultsView.vue` (post search)
    - `/favs` → `FavoritesView.vue` (user favorites)
  - Route guards handle query parameter parsing and state synchronization with the store
  - Integration with store's `searchPosts()` method for data loading on route change
- **Patterns**:
  - **Centralized Store**:
    - Single source of truth for application state
    - Uses reactive objects and computed properties for derived data
    - Synchronizes with backend via `/api/` endpoints
    - Handles caching (`posts`, `cachedTags`) to avoid redundant API calls
  - **Route-Driven State Updates**:
    - `router.beforeEach()` updates the store's search query and posts when navigating to search routes
    - Route parameters (`page`, `query`) are parsed and used to initialize store state
  - **Persistent Storage**:
    - `localStorage` is used for:
      - User authentication (`localStorage.getItem("account")`)
      - Settings (`saveSettings()`, `loadSettings()`)
      - Search history and favorites
  - **Error Handling**:
    - Toast notifications for errors (`store.toast`)
    - Error handling in API calls with `try/catch` and `fetch` rejection
## Key Files
- **Frontend**:
  - `client/src/store.ts` – Centralized reactive store with:
    - User state (`account`)
    - Search state (`query`, `posts`)
    - UI state (`settings`, `toast`)
    - API integration methods (`login()`, `searchPosts()`)
  - `client/src/router.ts` – Vue Router configuration with:
    - Route definitions and guards
    - Integration with store for route-based state updates
  - `client/src/views/*` – Vue SFCs for individual pages/components
## Code Sample: Store Integration
```ts
// Example from store.ts
export type SearchQuery = {
  include: Tag[];
  exclude: Tag[];
};
export type Store = {
  query: SearchQuery;
  searchPosts({ page, force }: { page?: number; force?: boolean; }): Promise<void>;
};
// Route integration in router.ts
router.beforeEach(to => {
  if (to.name === "search") {
    const page = parseInt(to.params.page as string, 10);
    const query = to.params.query as string[];
    tagsToSearchQuery(query).then(q => {
      store.query = q;
      store.searchPosts({ page, force: store.justClickedSearchButton });
    });
  }
});
```
