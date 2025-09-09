import { createRouter, createWebHistory } from "vue-router";

import FavoritesView from "./views/FavoritesView.vue";
import LandingView from "./views/LandingView.vue";
import SearchResultsView from "./views/SearchResultsView.vue";
import { tagsToSearchQuery } from "./search";
import store from "./store";

export const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: "/", name: "landing", component: LandingView },
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
    let shouldSearch = true;

    if (justLoaded) {
        shouldSearch = store.shouldSearchOnPageLoad();
        justLoaded = false;

        if (shouldSearch && to.name === "landing") {
            return router.push({
                name: "search",
                params: { page: 1 },
            });
        }
    }

    if (shouldSearch && to.name === "search") {
        const page = parseInt((to.params.page as string) || "1");
        const query = to.params.query || "";

        return new Promise<void>((resolve, reject) => {
            tagsToSearchQuery(query || []).then((q) => {
                store.query = q;
                store
                    .searchPosts(page, true)
                    .then(() => {
                        store.lastSearchRoute = to;
                        resolve();
                    })
                    .catch(reject);
            });
        });
    }
});
