import { createRouter, createWebHistory } from "vue-router";

import FavoritesView from "./views/FavoritesView.vue";
import LandingView from "./views/LandingView.vue";
import SearchResultsView from "./views/SearchResultsView.vue";
import { tagsToSearchQuery } from "./search";
import store from "./store";

export const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: "/",
            name: "landing",
            component: LandingView,
            beforeEnter(to) {
                if (!to.hash) {
                    return;
                }

                const queryParams = new URLSearchParams(
                    window.location.hash.replace(/^#/, ""),
                );

                // Redirect old URL scheme: /#page=N&q=tag,tag
                if (queryParams.has("page") && queryParams.has("q")) {
                    return {
                        name: "search",
                        params: {
                            page: queryParams.get("page"),
                            query: queryParams.get("q"),
                        },
                    };
                }
            },
        },
        {
            path: "/search/:page(\\d+)/:query?",
            name: "search",
            component: SearchResultsView,
        },
        { path: "/favs", name: "favorites", component: FavoritesView },
    ],
});

let justLoaded = true;

router.beforeEach((to, from) => {
    const alwaysSearchOnLoad =
        justLoaded &&
        store.settings.searchOnPageLoad === "always" &&
        to.name !== "favorites";
    const searchIfNonEmptyQuery =
        justLoaded &&
        to.name === "search" &&
        store.settings.searchOnPageLoad === "if-query" &&
        !!to.params.query;
    const navigatingToSearchPage = !justLoaded && to.name === "search";

    justLoaded = false;

    // Redirect to search page if we always search on load
    if (alwaysSearchOnLoad && to.name === "landing") {
        return router.push({
            name: "search",
            params: { page: 1 },
        });
    }

    const justClickedSearchButton = store.justClickedSearchButton;
    store.justClickedSearchButton = false;

    if (
        justClickedSearchButton ||
        alwaysSearchOnLoad ||
        searchIfNonEmptyQuery ||
        navigatingToSearchPage
    ) {
        const page = parseInt(to.params.page as string);

        return new Promise<void>((resolve, reject) => {
            tagsToSearchQuery(to.params.query || []).then((q) => {
                store.query = q;
                store
                    .searchPosts({ page, force: justClickedSearchButton })
                    .then(() => {
                        store.lastSearchRoute = to;
                        resolve();
                    })
                    .catch(reject);
            });
        });
    }
});
