import { createApp } from "vue";
import App from "./App.vue";
import { COMMIT_SHA } from "./config";
import { router } from "./router";
import store from "./store";
import { createRouterScroller } from "./vue-router-better-scroller";

// Periodically check the API and notify the user if the version updated
let currentVersion = COMMIT_SHA;
setInterval(() => {
	if (!store.settings.checkForUpdates) {
		return;
	}

	fetch("/api/version").then((resp) =>
		resp.json().then(({ version }: { version: string }) => {
			if (version !== currentVersion) {
				store.toast = {
					msg: "booruview updated, refresh the page",
					type: "info",
				};
				currentVersion = version;
			}
		}),
	);
}, 60 * 1000);

store.loadSettings();
store.updateCDNHosts();

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
