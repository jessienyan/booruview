<script setup lang="ts">
import store from "@/store";
import SidebarFooter from "./footer/SidebarFooter.vue";
import SidebarMain from "./main/SidebarMain.vue";
import { useNewFeatureIndicator } from "@/composable";
import { useTemplateRef } from "vue";
import { arrow, flip, offset, shift, useFloating } from "@floating-ui/vue";
import NewFeature from "../NewFeature.vue";

const favoritesButtonRef = useTemplateRef("favorites-btn");
const tooltipRef = useTemplateRef("tooltip");
const tooltipArrowRef = useTemplateRef("tooltip-arrow");
const { floatingStyles, middlewareData } = useFloating(
    favoritesButtonRef,
    tooltipRef,
    {
        placement: "right",
        middleware: [
            flip(),
            shift(),
            offset(10),
            arrow({ element: tooltipArrowRef }),
        ],
    },
);
const featFavorites = useNewFeatureIndicator(
    "favorites",
    new Date("2025-08-03"),
);
</script>

<template>
    <header class="sidebar-container" :class="{ closed: store.sidebarClosed }">
        <div class="buttons">
            <button
                class="btn-sidebar btn-toggle"
                @click="store.sidebarClosed = !store.sidebarClosed"
                title="toggle sidebar"
            >
                <i class="bi bi-list"></i>
            </button>

            <div class="spacing"></div>

            <button
                class="btn-sidebar btn-small"
                :class="{ active: store.postsBeingViewed === 'search-results' }"
                @click="store.postsBeingViewed = 'search-results'"
                title="view search results"
            >
                <i class="bi bi-search"></i>
            </button>
            <button
                class="btn-sidebar btn-small"
                :class="{ active: store.postsBeingViewed === 'favorites' }"
                @click="store.postsBeingViewed = 'favorites'"
                title="view favorites"
                ref="favorites-btn"
            >
                <i class="bi bi-heart"></i>
            </button>

            <div
                class="tooltip"
                v-if="featFavorites.show.value && store.sidebarClosed"
                ref="tooltip"
                :style="floatingStyles"
            >
                <p><NewFeature /></p>
                <p>you can find your favorites here</p>
                <p>
                    <button
                        class="btn-primary btn-rounded"
                        @click="featFavorites.onSeen()"
                    >
                        ok
                    </button>
                </p>

                <div
                    class="tooltip-arrow"
                    ref="tooltip-arrow"
                    :style="{
                        position: 'absolute',
                        left: '-30px',
                        top:
                            middlewareData.arrow?.y != null
                                ? `${middlewareData.arrow.y}px`
                                : '',
                    }"
                ></div>
            </div>
        </div>

        <div class="content" v-show="!store.sidebarClosed">
            <SidebarMain />
            <SidebarFooter />
        </div>
    </header>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/breakpoints";
@import "@/assets/colors";

.sidebar-container {
    display: flex;
    flex-direction: column;
    background-color: $color-sidebar;
    height: 100%;

    @media (max-width: $mobile-width) {
        .sidebar-open & {
            width: 100%;
        }

        .sidebar-closed & {
            height: auto;
        }
    }
}

.content {
    width: 450px;
    margin-top: 0.8rem;
    position: relative;
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;

    @media (max-width: $mobile-width) {
        .sidebar-open & {
            width: 100%;
        }
    }
}

$sidebar-width: 60px;

.buttons {
    display: flex;
    flex-direction: row;
    width: 100%;
    font-size: 48px;
    padding-right: 10px;
}

.btn-sidebar {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 48px;
    padding: 0;
    width: $sidebar-width;
    height: $sidebar-width;
    color: $color-primary-lighter;
}

.btn-small {
    font-size: 30px;
    width: $sidebar-width - 10px;

    &:hover,
    &.active {
        color: $color-primary-light;
    }

    &.active {
        text-shadow: 0 0 5px $color-primary-light;
    }
}

@media (min-width: $mobile-width) {
    .closed {
        .buttons {
            flex-direction: column;
            height: 100%;
            width: $sidebar-width;
        }

        .btn-small {
            width: $sidebar-width;
        }
    }
}

.spacing {
    margin: auto;
}

.tooltip {
    background-color: $color-primary;
    border: 1px solid $color-primary-lighter;
    border-radius: 5px;
    width: fit-content;
    padding: 0 1rem;
    filter: drop-shadow(0 0 10px black);
    z-index: 2;

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
    height: $arrowSize * 2;
    border: $arrowSize solid transparent;
    border-right: $arrowSize solid $color-primary;
}
</style>
