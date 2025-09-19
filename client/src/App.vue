<script setup lang="ts">
import store from "@/store";
import Sidebar from "@/components/sidebar/Sidebar.vue";
import { RouterView } from "vue-router";

import {
    computed,
    nextTick,
    provide,
    readonly,
    useTemplateRef,
    watch,
} from "vue";
import FullscreenView from "./components/fullscreen-view/FullscreenView.vue";
import ContentWarning from "./components/ContentWarning.vue";
import Toast from "./components/Toast.vue";
import { useNewFeatureIndicator } from "./composable";

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

const maintenanceBanner = useNewFeatureIndicator(
    "maintenance-banner-2025-09-19",
);
const maintenanceStart = new Date("2025-09-19T16:20:02.441Z");
const formattedDate = new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
}).format(maintenanceStart);
</script>

<template>
    <div class="app-outer">
        <div
            class="banner"
            v-if="store.settings.consented && maintenanceBanner.show.value"
        >
            Maintenance started @ {{ formattedDate }}. Booruview might be
            unavailable for a bit.
            <button @click="maintenanceBanner.onSeen()">
                <i class="bi bi-x-lg"></i> close
            </button>
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
                    'hide-overflow': store.fetchingPosts,
                }"
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
    height: 100%;
}

.banner {
    background-color: $color-primary;
    color: $color-primary-light;
    text-align: center;
    padding: 10px 20px;
    box-shadow: 0 0 20px black;

    button {
        border: 2px solid $color-primary-light;
        color: $color-primary-light;
        background: none;
        border-radius: 3px;
        font-weight: bold;
        cursor: pointer;
    }
}

.app {
    display: flex;
    flex-direction: row;
    flex: 1;
    overflow-y: scroll;

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
