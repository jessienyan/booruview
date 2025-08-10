<script setup lang="ts">
import store from "@/store";
import SidebarFooter from "./footer/SidebarFooter.vue";
import SidebarMain from "./main/SidebarMain.vue";
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
</style>
