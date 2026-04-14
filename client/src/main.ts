import { createApp } from "vue";
import App from "./App.vue";
import { COMMIT_SHA } from "./config";
import { router } from "./router";
import store, { APP_VERSION_TTL_MS } from "./store";
import { createRouterScroller } from "./vue-router-better-scroller";

// Periodically check the API and notify the user if the version updated
let currentVersion = COMMIT_SHA;
setInterval(async () => {
	if (!store.settings.checkForUpdates) {
		return;
	}

	try {
		const version = await store.appVersion();
		if(version !== currentVersion) {
			store.toast = {
				msg: "booruview updated, refresh the page",
				type: "info",
			};
			currentVersion = version;
		}
	} catch (_e) {}
}, APP_VERSION_TTL_MS);

store.loadSettings();
store.updateCDNHosts();
store.fetchAccountData();

const app = createApp(App);
app.use(router);
app.use(
	createRouterScroller({
		selectors: {
			"#scroll-container": ({ element, savedPosition }) => {
				savedPosition = savedPosition ?? { top: 0 };

				// Try scrolling immediately
				element.scrollTo(savedPosition);

				// Scroll again after a short delay in case the component is slow to re-render.
				// The delayed scroll on its own can be a bit disorienting
				setTimeout(() => element.scrollTo(savedPosition), 20);
			},
		},
	}),
);
app.mount(document.body);
