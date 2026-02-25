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
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="currentColor"
                    class="bi"
                    viewBox="0 0 16 16"
                >
                    <path
                        fill-rule="evenodd"
                        d="M2.5 12a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5"
                    />
                </svg>
            </button>

            <div class="spacing"></div>

            <RouterLink :to="store.lastSearchRoute || { name: 'landing' }">
                <button class="btn-sidebar btn-nav" title="view search results">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        class="bi"
                        viewBox="0 0 16 16"
                    >
                        <path
                            d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"
                        />
                    </svg>
                </button>
            </RouterLink>
            <RouterLink :to="{ name: 'favorites' }">
                <button
                    class="btn-sidebar btn-nav"
                    title="view favorites"
                    ref="favorites-btn"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        class="bi"
                        viewBox="0 0 16 16"
                    >
                        <path
                            d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.01zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143q.09.083.176.171a3 3 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15"
                        />
                    </svg>
                </button>
            </RouterLink>
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

$sidebar-btn-width: 60px;

.buttons {
    display: flex;
    flex-direction: row;
    width: 100%;
    font-size: 0;
    padding-right: 10px;

    .router-link-active {
        .btn-sidebar {
            color: $color-primary-light;
        }

        .btn-sidebar {
            text-shadow: 0 0 5px $color-primary-light;
        }
    }
}

.btn-sidebar {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 0;
    width: 60px;
    height: 60px;
    color: $color-primary-lighter;

    .bi {
        width: 100%;
        height: 100%;
    }
}

.btn-toggle {
    padding: 6px;
}

.btn-nav {
    padding: 14px 9px;
    width: $sidebar-btn-width - 10px;
}

@media (min-width: $mobile-width) {
    .closed {
        .buttons {
            flex-direction: column;
            height: 100%;
            width: $sidebar-btn-width;
        }

        .btn-small {
            width: $sidebar-btn-width;
        }
    }
}

.spacing {
    margin: auto;
}

.survey {
    background-color: $color-primary;
    padding: 15px;
    margin-bottom: 15px;

    p {
        color: $color-primary-light;

        &:first-child {
            margin-top: 0;
        }

        &:last-child {
            margin-bottom: 0;
        }
    }

    .hide-survey {
        float: right;
        border: none;
        background: none;
        color: $color-primary-light;
    }
}
</style>
