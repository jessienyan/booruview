import { createApp } from "vue";
import App from "./App.vue";
import store from "./store";
import { COMMIT_SHA } from "./config";

// Periodically check the API and notify the user if the version updated
let currentVersion = COMMIT_SHA;
setInterval(() => {
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

const app = createApp(App);
app.mount(document.body);
