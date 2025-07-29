<script setup lang="ts">
import {
    computed,
    onMounted,
    onUnmounted,
    ref,
    useTemplateRef,
    watch,
} from "vue";
import store from "@/store";
import PostContainer from "./components/PostContainer.vue";
import Sidebar from "./components/sidebar/Sidebar.vue";
import FullscreenView from "./components/fullscreen-view/FullscreenView.vue";
import Router from "./Router.vue";
import ContentWarning from "./components/ContentWarning.vue";
import Footer from "./components/Footer.vue";
import NoResults from "./components/NoResults.vue";
import ChipMenu from "./components/tag-chip/ChipMenu.vue";
import Toast from "./components/Toast.vue";
import PageChangeGesture from "./PageChangeGesture.vue";

const mainContainer = useTemplateRef("main");
const scrollPositionHistory = ref<{ [page: number]: number }>({});

function onPostsCleared() {
    // Reset any scroll position history when a new search is triggered
    scrollPositionHistory.value = {};
}

onMounted(() => {
    store.onPostsCleared.addEventListener("postsCleared", onPostsCleared);
});

onUnmounted(() => {
    store.onPostsCleared.removeEventListener("postsCleared", onPostsCleared);
});

// Pre-DOM update on page change
watch(
    () => store.currentPage,
    (_, prevPage) => {
        // Remember where the scroll position was when the page changes. This watcher
        // runs before the DOM has updated
        scrollPositionHistory.value[prevPage] =
            mainContainer.value?.scrollTop ?? 0;
    },
);

// Post-DOM update on page change
watch(
    () => [store.currentPage, store.posts.size],
    () => {
        // Prevent firing on page load if there are no posts yet
        if (store.posts.size === 0) {
            return;
        }

        // Scroll needs to be deferred in order to work on mobile
        setTimeout(() => {
            if (mainContainer.value === null) {
                return;
            }

            // If this is a new page scroll to the top
            if (!(store.currentPage in scrollPositionHistory.value)) {
                scrollPositionHistory.value[store.currentPage] = 0;
            }

            // Restore scroll position
            mainContainer.value.scrollTop =
                scrollPositionHistory.value[store.currentPage];
        }, 0);

        mainContainer.value!.focus();
    },
    {
        flush: "post",
    },
);

// Focus main container when exiting fullscreen view
watch(
    () => store.fullscreenPost,
    () => {
        if (store.fullscreenPost === null) {
            mainContainer.value!.focus();
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
    <Router />
    <PageChangeGesture
        v-if="mainContainer != null"
        :scroll-container="mainContainer!"
    />
    <div
        class="app"
        :class="{
            'sidebar-closed': store.sidebarClosed,
            'sidebar-open': !store.sidebarClosed,
        }"
    >
        <Toast
            v-if="store.toast.msg.length > 0"
            :kind="store.toast.type"
            @dismiss="store.toast.msg = ''"
            >{{ store.toast.msg }}</Toast
        >
        <ChipMenu v-if="store.tagMenu !== null" />
        <ContentWarning v-if="!hasConsented" />

        <FullscreenView
            v-if="store.fullscreenPost !== null"
            :post="store.fullscreenPost"
        />

        <Sidebar
            :closed="store.sidebarClosed"
            @toggle="store.sidebarClosed = !store.sidebarClosed"
        />
        <main ref="main" tabindex="-1">
            <template v-if="store.hasSearched">
                <NoResults v-if="store.totalPostCount === 0" />
                <template v-else>
                    <PostContainer
                        :posts="store.postsForCurrentPage() || []"
                        :scroll-container="mainContainer!"
                    />
                    <Footer v-if="!store.fetchingPosts" />
                </template>
            </template>
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

    @media (max-width: $mobile-width) {
        .sidebar-open & {
            display: none;
        }
    }

    &:focus {
        outline: none;
    }
}
</style>
