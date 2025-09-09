import { createRouter, createWebHashHistory } from "vue-router";

import FavoritesView from "./views/FavoritesView.vue";
import LandingView from "./views/LandingView.vue";
import SearchResultsView from "./views/SearchResultsView.vue";
import store from "./store";

export const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        { path: "/", name: "landing", component: LandingView },
        {
            path: "/search/:page(\\d+)/:query?",
            name: "search",
            component: SearchResultsView,
            props: true,
            beforeEnter(to, from) {
                store.lastSearchRoute = to;
            },
        },
        { path: "/favs", name: "favorites", component: FavoritesView },
    ],
});
