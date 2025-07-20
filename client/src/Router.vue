<script setup lang="ts">
import { onMounted } from "vue";
import store from "./store";

function parsePage(raw: string): number {
    const val = parseInt(raw);
    if (!Number.isSafeInteger(val) || val < 1) {
        return 1;
    }

    return val;
}

function parseQuery(raw: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
        if (raw.length === 0) {
            store.query.clear();
            resolve();
            return;
        }

        const tagNames = raw.split(",").filter((v) => v.length > 0);

        store
            .loadTags(tagNames)
            .then(() => {
                store.query.clear();

                for (let name of tagNames) {
                    const negate = name[0] === "-";
                    if (negate) {
                        name = name.slice(1);
                    }

                    let tag = store.cachedTags.get(name);
                    if (tag === undefined) {
                        tag = {
                            count: 0,
                            name: name,
                            type: "unknown",
                        };
                    }

                    if (negate) {
                        store.query.excludeTag(tag);
                    } else {
                        store.query.includeTag(tag);
                    }
                }

                resolve();
            })
            .catch(reject);
    });
}

// Loads the query params into the store
function loadQueryParams(): Promise<void> {
    return new Promise((resolve, reject) => {
        const params = new URLSearchParams(
            window.location.hash.replace(/^#/, ""),
        );
        const page = params.get("page") || "1";
        const query = params.get("q") || "";
        store.currentPage = parsePage(page);
        parseQuery(query).then(resolve).catch(reject);
    });
}

function onRouteChange() {
    const lastPage = store.currentPage;
    const lastSearch = store.query.copy();

    loadQueryParams()
        .then(() => {
            // NOTE: Checking if these changed can help avoid double searches. For example when
            // clicking the "next" button while viewing a post in fullscreen mode causes a page change
            const pageChanged = store.currentPage !== lastPage;
            const searchChanged = !lastSearch.equals(store.query);

            if (searchChanged) {
                store.posts.clear();
            }

            if (pageChanged || searchChanged) {
                // Fix fullscreen not being dismissed when using browser navigation
                store.fullscreenPost = null;
                store.searchPosts();
            }
        })
        .catch(() => {});
}

function onPageLoad() {
    loadQueryParams()
        .then(() => {
            if (store.settings.consented && store.shouldSearchOnPageLoad()) {
                store.searchPosts();
            }
        })
        .catch(() => {
            store.setQueryParams();
        });
}

window.addEventListener("hashchange", onRouteChange);
onMounted(onPageLoad);
</script>

<template></template>
