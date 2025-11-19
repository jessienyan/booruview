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

router.beforeEach((to) => {
	if (to.name === "search") {
		if (!store.settings.consented) {
			return;
		}

		const page = parseInt(to.params.page as string);

		return new Promise<void>((resolve, reject) => {
			tagsToSearchQuery(to.params.query || []).then((q) => {
				store.query = q;
				store
					.searchPosts({ page, force: store.justClickedSearchButton })
					.then(() => {
						store.lastSearchRoute = to;
						resolve();
					})
					.catch(reject)
					.finally(() => (store.justClickedSearchButton = false));
			});
		});
	}
});
