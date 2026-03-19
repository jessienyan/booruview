import { createRouter, createWebHistory } from "vue-router";
import { tagsToSearchQuery } from "./search";
import store from "./store";
import FavoritesView from "./views/FavoritesView.vue";
import LandingView from "./views/LandingView.vue";
import SearchResultsView from "./views/SearchResultsView.vue";

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: "/",
			name: "landing",
			component: LandingView,
		},
		{
			path: "/search/:page(\\d+)/:query?",
			name: "search",
			component: SearchResultsView,
		},
		{
			path: "/favs/:page(\\d+)?",
			name: "favorites",
			component: FavoritesView,
		},
	],
});

router.beforeEach(async to => {
	if (to.name === "search") {
		if (!store.settings.consented) {
			return;
		}

		const page = parseInt(to.params.page as string, 10);
		let query: string[];

		if(!to.params.query || to.params.query.length === 0) {
			query = [];
		} else if(!Array.isArray(to.params.query)) {
			query = to.params.query.split(",");
		} else {
			query = to.params.query;
		}

		try {
			const q = await tagsToSearchQuery(query || []);
			store.query = q;
			await store.searchPosts({ page, force: store.justClickedSearchButton });
			store.lastSearchRoute = to;
		} finally {
			store.justClickedSearchButton = false;
		}
	}
});
