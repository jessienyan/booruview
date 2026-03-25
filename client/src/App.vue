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
</script>

<template>
    <div class="app-outer">
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

            <Transition name="transition">
                <ContentWarning v-if="!hasConsented" />
            </Transition>

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

.banner {
    padding: 1em;
    font-size: 1.1em;
    line-height: 1.5em;
}

.banner-warning {
    background-color: #ccd666;
    color: #111;
}

.banner-announcement {
    background-color: $color-primary-light;
    color: #111;
}

.banner-close {
    display: inline-block;
    padding: 4px;
    font-weight: bold;
    text-decoration: underline;
    color: inherit;
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
