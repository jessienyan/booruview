import { createApp } from "vue";
import App from "./App.vue";
import store from "./store";
import { COMMIT_SHA } from "./config";
import { router } from "./router";
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

const app = createApp(App);
app.use(router);
app.use(
    createRouterScroller({
        selectors: {
            "#scroll-container": ({ element, savedPosition }) => {
                // Add a short delay in case the container needs to re-render a couple times
                setTimeout(
                    () => element.scrollTo(savedPosition ?? { top: 0 }),
                    20,
                );
            },
        },
    }),
);
app.mount(document.body);
