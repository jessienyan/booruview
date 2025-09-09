import { createRouter, createWebHashHistory } from "vue-router";

import LandingView from "./views/LandingView.vue";
import SearchResultsView from "./views/SearchResultsView.vue";

export const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        { path: "/", component: LandingView },
        { path: "/search/:page", component: SearchResultsView, props: true },
        {
            path: "/search/:page/:query",
            component: SearchResultsView,
            props: true,
        },
    ],
});
