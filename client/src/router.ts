import { createRouter, createWebHistory } from "vue-router";
import { POSTS_PER_PAGE } from "./config";
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

			await store.searchAndUpdateResults({ query, page, force: store.justClickedSearchButton });
			store.lastSearchRoute = to;
		} finally {
			store.justClickedSearchButton = false;
		}
	}

	if (to.name === "favorites") {
		const favPosts = store.favoritePosts();
		const maxPage = Math.ceil(favPosts.value.length / POSTS_PER_PAGE);

		if (maxPage === 0) {
			if (to.params.page) {
				return { name: "favorites" };
			}
			return;
		}

		let page = parseInt(to.params.page as string || "1", 10);
		page = Math.max(1, Math.min(page, maxPage));
		store.lastFavPage = page;

		if (page.toString() !== to.params.page) {
			return {
				name: "favorites",
				params: { page },
			};
		}
	}
});
