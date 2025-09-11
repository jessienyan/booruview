<script setup lang="ts">
import store from "@/store";
import Sidebar from "@/components/sidebar/Sidebar.vue";
import { RouterView } from "vue-router";

import { computed, provide, readonly, useTemplateRef } from "vue";
import FullscreenView from "./components/fullscreen-view/FullscreenView.vue";
import ContentWarning from "./components/ContentWarning.vue";
import ChipMenu from "./components/tag-chip/ChipMenu.vue";
import Toast from "./components/Toast.vue";

const mainContainer = useTemplateRef("main");
provide("mainContainer", readonly(mainContainer));

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
        <ChipMenu v-if="store.tagMenu !== null" />
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
            }"
        >
            <RouterView />
        </main>
    </div>
</template>

<style scoped lang="scss">
@import "@/assets/breakpoints";
@import "@/assets/colors";

.app {
    display: flex;
    flex-direction: row;
    width: 100%;
    height: 100%;

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
}

.prevent-pull-to-refresh {
    overscroll-behavior-y: contain;
}
</style>
