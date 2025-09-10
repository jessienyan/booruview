/*

including this directly for now until the PR i submitted gets merged
https://github.com/antfu/vue-router-better-scroller/pull/5

*/

import { nextTick } from "vue";
import { isNavigationFailure } from "vue-router";

const STATE_KEY = "vueRouterScroller";
const DEFAULT_INTERVAL = 200;
function setupRouterScroller(router, options) {
    if (router.options.scrollBehavior)
        console.warn(
            "`scrollBehavior` options in Vue Router is overwritten by `vue-router-scroller` plugin, you can remove it from createRouter()",
        );
    router.options.scrollBehavior = () => {};
    options.storeInterval = options.storeInterval ?? DEFAULT_INTERVAL;
    const positionsMap = /* @__PURE__ */ new Map();
    let saveTimerId = null;
    function startTimer() {
        if (saveTimerId !== null) return;
        saveTimerId = window.setInterval(() => {
            const pos = capturePositions(options);
            const key = getScrollKey(history.state.fullPath);
            positionsMap.set(key, pos);
            history.replaceState({ ...history.state, [STATE_KEY]: pos }, "");
        }, options.storeInterval);
    }
    function stopTimer() {
        if (saveTimerId === null) return;
        clearInterval(saveTimerId);
        saveTimerId = null;
    }
    startTimer();
    window.addEventListener("popstate", () => stopTimer());
    router.beforeResolve(() => stopTimer());
    router.afterEach((to, from, failure) => {
        if (isNavigationFailure(failure)) return;
        const key = getScrollKey(to.fullPath);
        const pos = history.state[STATE_KEY] || positionsMap.get(key);
        const type = history.state[STATE_KEY] ? "history" : "push";
        nextTick(() => {
            applyPositions(to, from, pos, type, options);
            startTimer();
        });
    });
}
function createRouterScroller(options) {
    return {
        install(app) {
            const router = app.config.globalProperties.$router;
            if (!router)
                throw new Error(
                    "Router instance is not found on this Vue app. This plugin should be installed after Vue Router.",
                );
            setupRouterScroller(router, options);
        },
    };
}
function getScrollKey(path, delta = 0) {
    const position = history.state ? history.state.position - delta : -1;
    return position + path;
}
function capturePositions(options) {
    const pos = {};
    for (const [selector] of Object.entries(options.selectors)) {
        const element = querySelector(selector);
        if (!element) continue;
        pos[selector] = getScrollPosition(element);
    }
    return pos;
}
function querySelector(name) {
    if (typeof window === "undefined") return void 0;
    if (name === "body") return document.body;
    if (name === "window") return window;
    return document.querySelector(name);
}
function getScrollPosition(el) {
    if (el instanceof Window)
        return { left: window.scrollX, top: window.scrollY };
    else return { left: el.scrollLeft, top: el.scrollTop };
}
async function applyPositions(to, from, pos, type, options) {
    for (const [selector, handler] of Object.entries(options.selectors)) {
        const element = querySelector(selector);
        if (!element) continue;
        let position = pos?.[selector];
        if (typeof handler === "function") {
            const result = await handler({
                to,
                from,
                element,
                selector,
                type,
                savedPosition: position,
            });
            if (!result) continue;
            if (result !== true) position = result;
        } else if (handler === true) {
            if (type === "push") position = void 0;
        }
        element.scrollTo({
            behavior: options.behavior,
            ...(position || { top: 0, left: 0 }),
        });
    }
}

export { createRouterScroller, setupRouterScroller };
