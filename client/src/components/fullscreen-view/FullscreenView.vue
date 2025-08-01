<script setup lang="ts">
import store, { type FullscreenViewMenuAnchorPoint } from "@/store";
import {
    computed,
    onMounted,
    onUnmounted,
    ref,
    useTemplateRef,
    type CSSProperties,
} from "vue";
import ContentTab from "./ContentTab.vue";
import InfoTab from "./InfoTab.vue";
import ScreenCover from "../ScreenCover.vue";
import {
    useIsVideo,
    useNewFeatureIndicator,
    useStationaryClick,
} from "@/composable";
import { arrow, flip, offset, shift, useFloating } from "@floating-ui/vue";
import NewFeature from "../NewFeature.vue";

type Tab = "content" | "info";
type PostNavInfo = { page: number; index: number } | null;

const { post } = defineProps<{ post: Post }>();
const currentTab = ref<Tab>("content");

const borderRadius = "15px";

const menuAnchorPoints: Record<FullscreenViewMenuAnchorPoint, CSSProperties> = {
    topleft: { top: "0px", left: "0px", borderBottomRightRadius: borderRadius },
    topcenter: {
        top: "0px",
        left: "50%",
        transform: "translateX(-50%)",
        borderBottomLeftRadius: borderRadius,
        borderBottomRightRadius: borderRadius,
    },
    topright: {
        top: "0px",
        right: "0px",
        borderBottomLeftRadius: borderRadius,
    },
    right: {
        top: "50%",
        right: "0px",
        transform: "translateY(-50%)",
        borderTopLeftRadius: borderRadius,
        borderBottomLeftRadius: borderRadius,
    },
    bottomright: {
        bottom: "0px",
        right: "0px",
        borderTopLeftRadius: borderRadius,
    },
    bottomcenter: {
        bottom: "0px",
        left: "50%",
        transform: "translateX(-50%)",
        borderTopLeftRadius: borderRadius,
        borderTopRightRadius: borderRadius,
    },
    bottomleft: {
        bottom: "0px",
        left: "0px",
        borderTopRightRadius: borderRadius,
    },
    left: {
        top: "50%",
        left: "0px",
        transform: "translateY(-50%)",
        borderTopRightRadius: borderRadius,
        borderBottomRightRadius: borderRadius,
    },
};

const featMenuAnchor = useNewFeatureIndicator(
    "menu-anchor",
    new Date("2025-07-27"),
);
const menuRef = useTemplateRef("menu");
const tooltipRef = useTemplateRef("tooltip");
const tooltipArrowRef = useTemplateRef("tooltip-arrow");
const { floatingStyles, middlewareData } = useFloating(menuRef, tooltipRef, {
    middleware: [
        flip(),
        shift(),
        offset(20),
        arrow({ element: tooltipArrowRef }),
    ],
});

function close() {
    store.fullscreenPost = null;
    // Explicitly close the tag menu since panzoom eats the click event
    store.tagMenu = null;
}

function onKeyDown(e: KeyboardEvent) {
    if (e.key === "Esc" || e.key === "Escape") {
        e.preventDefault();
        close();
    } else if (e.key === "ArrowLeft" || e.key.toUpperCase() === "A") {
        e.preventDefault();
        showPrevPost();
    } else if (e.key === "ArrowRight" || e.key.toUpperCase() === "D") {
        e.preventDefault();
        showNextPost();
    } else if (e.key.toUpperCase() === "F") {
        e.preventDefault();
        toggleFavorite();
    }
}

const isVideo = useIsVideo(() => post);

const tabClasses = computed(() => {
    // Center unless we're displaying an image and viewing the content tab.
    // Image viewing uses panzoom which breaks when the container is flexbox
    const centered = !(currentTab.value === "content" && !isVideo.value);
    return {
        "tab-centered": centered,
    };
});

const tabHandler = useStationaryClick(close);

const currentPostIndex = computed(() => {
    let index = -1;

    if (store.postsBeingViewed === "search-results") {
        index = store.posts
            .get(store.currentPage)!
            .findIndex((p) => p.id === post.id);
    } else if (store.postsBeingViewed === "favorites") {
        index = store.settings.favorites.findIndex((p) => p.id === post.id);
    }

    return index === -1 ? null : index;
});
const nextPost = computed<PostNavInfo>(() => {
    if (currentPostIndex.value == null) {
        return null;
    }

    if (store.postsBeingViewed === "favorites") {
        if (currentPostIndex.value === store.settings.favorites.length - 1) {
            return null;
        }

        return {
            page: 1,
            index: currentPostIndex.value + 1,
        };
    }

    const isLastPage = store.currentPage === store.maxPage();
    const isLastResult = currentPostIndex.value === store.resultsPerPage - 1;

    if (isLastPage && isLastResult) {
        return null;
    }

    if (isLastResult) {
        return {
            page: store.currentPage + 1,
            index: 0,
        };
    }

    return {
        page: store.currentPage,
        index: currentPostIndex.value + 1,
    };
});

const prevPost = computed<PostNavInfo>(() => {
    if (currentPostIndex.value == null) {
        return null;
    }

    if (store.postsBeingViewed === "favorites") {
        if (currentPostIndex.value === 0) {
            return null;
        }

        return {
            page: 1,
            index: currentPostIndex.value - 1,
        };
    }

    const isFirstPage = store.currentPage === 1;
    const isFirstResult = currentPostIndex.value === 0;

    if (isFirstPage && isFirstResult) {
        return null;
    }

    if (isFirstResult) {
        return {
            page: store.currentPage - 1,
            index: store.resultsPerPage - 1,
        };
    }

    return {
        page: store.currentPage,
        index: currentPostIndex.value - 1,
    };
});

function showNextPost() {
    const nav = nextPost.value;

    if (nav === null) {
        return;
    } else if (store.postsBeingViewed === "favorites") {
        store.fullscreenPost = store.settings.favorites[nav.index];
        return;
    }

    if (store.currentPage === nav.page) {
        store.fullscreenPost = store.posts.get(nav.page)![nav.index] || null;
    } else {
        store
            .nextPage()
            ?.then(
                () =>
                    (store.fullscreenPost =
                        store.posts.get(nav.page)![nav.index] || null),
            );
    }
}

function showPrevPost() {
    const nav = prevPost.value;

    if (nav === null) {
        return;
    } else if (store.postsBeingViewed === "favorites") {
        store.fullscreenPost = store.settings.favorites[nav.index];
        return;
    }

    if (store.currentPage === nav.page) {
        store.fullscreenPost = store.posts.get(nav.page)![nav.index] || null;
    } else {
        store
            .prevPage()
            ?.then(
                () =>
                    (store.fullscreenPost =
                        store.posts.get(nav.page)![nav.index] || null),
            );
    }
}

function toggleFavorite() {
    if (isFavorited.value) {
        store.settings.favorites.splice(favoriteIndex.value, 1);
    } else {
        store.settings.favorites = [post].concat(store.settings.favorites);
    }

    store.settings.save();
}

const favoriteIndex = computed(() =>
    store.settings.favorites.findIndex((p) => p.id === post.id),
);
const isFavorited = computed(() => favoriteIndex.value !== -1);

onMounted(() => {
    document.addEventListener("keydown", onKeyDown, { capture: true });
});

onUnmounted(() => {
    document.removeEventListener("keydown", onKeyDown, { capture: true });
});
</script>

<template>
    <div class="fullscreen-viewer">
        <ScreenCover />
        <div class="viewer-container">
            <div
                class="tab"
                :class="tabClasses"
                @mousedown.self="tabHandler.mouseDown"
                @mouseup.self="tabHandler.mouseUp"
            >
                <KeepAlive>
                    <ContentTab v-if="currentTab == 'content'" :post="post" />
                    <InfoTab v-else-if="currentTab == 'info'" :post="post" />
                </KeepAlive>
            </div>

            <div
                class="tooltip"
                v-if="featMenuAnchor.show.value"
                ref="tooltip"
                :style="floatingStyles"
            >
                <p><NewFeature /></p>
                <p>
                    this menu can now be moved and rotated. check the settings
                    tab in the sidebar
                </p>
                <p>
                    <button
                        class="btn-primary btn-rounded"
                        @click="featMenuAnchor.onSeen()"
                    >
                        ok
                    </button>
                </p>

                <div
                    class="tooltip-arrow"
                    ref="tooltip-arrow"
                    :style="{
                        position: 'absolute',
                        left:
                            middlewareData.arrow?.x != null
                                ? `${middlewareData.arrow.x}px`
                                : '',
                        top:
                            middlewareData.arrow?.y != null
                                ? `${middlewareData.arrow.y}px`
                                : '',
                        transform:
                            middlewareData.arrow?.y == null
                                ? 'rotate(180deg)'
                                : '',
                    }"
                ></div>
            </div>

            <footer
                ref="menu"
                class="tab-menu"
                :class="{ flipped: store.settings.fullscreenViewMenuRotate }"
                :style="
                    menuAnchorPoints[store.settings.fullscreenViewMenuAnchor]
                "
            >
                <button
                    class="menu-btn"
                    :class="{ active: currentTab == 'content' }"
                    @click="currentTab = 'content'"
                    title="view content"
                >
                    <i class="bi bi-image"></i>
                </button>
                <button
                    class="menu-btn"
                    :class="{ active: currentTab == 'info' }"
                    @click="currentTab = 'info'"
                    title="view tags"
                >
                    <i class="bi bi-info-circle"></i>
                </button>
                <button
                    class="menu-btn fav-btn"
                    :class="{ favorited: isFavorited }"
                    @click="toggleFavorite"
                    title="favorite"
                >
                    <i
                        class="bi"
                        :class="[`bi-${isFavorited ? 'heart-fill' : 'heart'}`]"
                    ></i>
                </button>
                <button
                    class="menu-btn"
                    title="previous image"
                    @click="showPrevPost"
                    :disabled="prevPost === null"
                >
                    <i class="bi bi-arrow-left"></i>
                </button>
                <button
                    class="menu-btn"
                    title="next image"
                    @click="showNextPost"
                    :disabled="nextPost === null"
                >
                    <i class="bi bi-arrow-right"></i>
                </button>
                <button
                    class="menu-btn close-btn"
                    @click="close()"
                    title="close"
                >
                    <i class="bi bi-x-lg"></i>
                </button>
            </footer>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/breakpoints";
@import "@/assets/buttons";
@import "@/assets/colors";

.fullscreen-viewer {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    overflow: hidden;
}

.viewer-container {
    height: 100%;
    position: relative;
    z-index: 2;
    width: 100%;
}

.tooltip {
    background-color: $color-primary;
    border: 1px solid $color-primary-lighter;
    border-radius: 5px;
    width: 300px;
    padding: 0 1rem;
    filter: drop-shadow(0 0 10px black);

    p {
        color: $color-primary-light;
    }

    button {
        width: 100%;
    }
}

.tooltip-arrow {
    $arrowSize: 15px;

    width: $arrowSize * 2;
    height: $arrowSize;
    border-left: $arrowSize solid transparent;
    border-right: $arrowSize solid transparent;
    border-bottom: $arrowSize solid $color-primary;
}

.tab {
    min-height: 0;
    height: 100%;
}

.tab-centered {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 4rem 0.8rem;
}

.tab-menu {
    display: flex;
    position: absolute;
    z-index: 2;
    background-color: rgba(0, 0, 0, 0.7);
    box-shadow: 0 0 0.8rem black;

    &.flipped {
        flex-direction: column-reverse;
    }
}

.menu-btn {
    background: none;
    border: none;
    color: white;
    font-size: 30px;
    text-shadow: 0 0 0.4rem white;
    opacity: 0.5;
    transition: opacity 150ms;
    cursor: pointer;
    padding: 0.6rem 0.8rem;

    &:not(:disabled) {
        &:hover,
        &.active {
            opacity: 1;
        }
    }

    &:disabled {
        opacity: 0.2;
        cursor: default;
    }
}

.fav-btn {
    .bi {
        position: relative;
        top: 2px;
    }

    &.favorited {
        color: pink;
    }
}

.close-btn {
    color: #bb9fce;
    text-shadow: 0 0 0.4rem #bb9fce;
    padding-right: 0.8rem;
}
</style>
