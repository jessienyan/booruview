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
			beforeEnter(to) {
				if (!to.hash) {
					return;
				}

				const queryParams = new URLSearchParams(window.location.hash.replace(/^#/, ""));

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

router.beforeEach(async (to, from) => {
	if (to.name === "search") {
		if (!store.settings.consented) {
			return;
		}

		const page = parseInt(to.params.page as string, 10);
		let rawQuery: string[];

		if(!to.params.query || to.params.query.length === 0) {
			rawQuery = [];
		} else if(!Array.isArray(to.params.query)) {
			rawQuery = to.params.query.split(",");
		} else {
			rawQuery = to.params.query;
		}

		try {
			const query = await tagsToSearchQuery(rawQuery || []);

			// Overwrite stored query only if it changed in the URL. This lets users edit
			// their search between pages or when pressing forward/back in the browser.
			// The query won't be applied until they click search
			if(to.params.query !== from.params.query) {
				store.query = query;
			}

			await store.searchPosts({ query, page, force: store.justClickedSearchButton });
			store.lastSearchRoute = to;
		} finally {
			store.justClickedSearchButton = false;
		}
	}
});
