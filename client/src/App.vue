<script setup lang="ts">
import {
    computed,
    nextTick,
    provide,
    readonly,
    useTemplateRef,
    watch,
} from "vue";
import { RouterView } from "vue-router";
import Sidebar from "@/components/sidebar/Sidebar.vue";
import { mediaProxyBanner } from "@/indicators";
import store from "@/store";
import ContentWarning from "./components/ContentWarning.vue";
import FullscreenView from "./components/fullscreen-view/FullscreenView.vue";
import Toast from "./components/Toast.vue";

const mainContainer = useTemplateRef("main");
provide("mainContainer", readonly(mainContainer));

// Focus scroll container when sidebar is closed
watch(
    () => store.sidebarClosed,
    () => {
        if (store.sidebarClosed) {
            nextTick(() => mainContainer.value?.focus());
        }
    },
);

const hasConsented = computed(() => {
    if (store.settings.consented) {
        return true;
    }

    // Don't show consent modal for search engine crawlers
    const crawlers = /Googlebot|Bingbot|DuckDuckbot/;
    const isCrawler = crawlers.exec(navigator.userAgent) !== null;
    return isCrawler;
});

const bannerDate = new Date("2026-02-16").toLocaleDateString();
</script>

<template>
    <div class="app-outer">
        <div
            class="banner-warning"
            v-if="hasConsented && mediaProxyBanner.show.value"
        >
            As of {{ bannerDate }}, Gelbooru blocked direct linking, so all
            media is now proxied through <code>proxy.booruview.com</code>. HD
            images are temporarily disabled to save on bandwidth. Report any
            issues in the "about" tab, and consider donating to help with
            hosting costs. Thank you ❤️
            <a
                href="#"
                class="banner-close"
                @click.prevent="mediaProxyBanner.onHide()"
                >close</a
            >
        </div>

        <div
            class="app"
            :class="{
                'sidebar-closed': store.sidebarClosed,
                'sidebar-open': !store.sidebarClosed,
            }"
        >
            <Transition>
                <Toast
                    v-if="store.toast.msg.length > 0"
                    :kind="store.toast.type"
                    @dismiss="store.toast.msg = ''"
                    >{{ store.toast.msg }}</Toast
                >
            </Transition>
            <ContentWarning v-if="!hasConsented" />

            <FullscreenView
                v-if="store.fullscreenPost !== null"
                :post="store.fullscreenPost"
            />
            <Sidebar />
            <main
                id="scroll-container"
                ref="main"
                tabindex="-1"
                :class="{
                    'prevent-pull-to-refresh': store.userIsSwipingToChangePage,
                    'hide-overflow':
                        store.fetchingPosts || store.userIsSwipingToChangePage,
                }"
                @keydown.left="$route.name === 'search' && store.prevPage()"
                @keydown.right="$route.name === 'search' && store.nextPage()"
            >
                <RouterView />
            </main>
        </div>
    </div>
</template>

<style scoped lang="scss">
@import "@/assets/breakpoints";
@import "@/assets/colors";

.app-outer {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
}

.banner-warning {
    background-color: #ccd666;
    color: #111;
    padding: 1em;
    font-size: 1.1em;
    line-height: 1.5em;

    .banner-close {
        font-weight: bold;
        text-decoration: underline;
        color: inherit;
    }
}

.app {
    display: flex;
    flex-direction: row;
    flex: 1;
    min-height: 0;

    @media (max-width: $mobile-width) {
        &.sidebar-closed {
            flex-direction: column;
        }
    }
}

main {
    flex: 1;
    min-height: 0;
    overflow-y: scroll;
    overscroll-behavior-x: contain;
    position: relative;

    @media (max-width: $mobile-width) {
        .sidebar-open & {
            display: none;
        }
    }

    &:focus {
        outline: none;
    }

    &.hide-overflow {
        overflow-y: hidden;
    }
}

.prevent-pull-to-refresh {
    overscroll-behavior-y: contain;
}
</style>
