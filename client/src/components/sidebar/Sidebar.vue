<script setup lang="ts">
import store from "@/store";
import SidebarFooter from "./footer/SidebarFooter.vue";
import SidebarMain from "./main/SidebarMain.vue";
import { useDontShowAgain } from "@/composable";
import { SURVEY_LINK } from "@/config";

const survey = useDontShowAgain("hide-survey");
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

            <RouterLink :to="store.lastSearchRoute || { name: 'landing' }">
                <button
                    class="btn-sidebar btn-small"
                    title="view search results"
                >
                    <i class="bi bi-search"></i>
                </button>
            </RouterLink>
            <RouterLink :to="{ name: 'favorites' }">
                <button
                    class="btn-sidebar btn-small"
                    title="view favorites"
                    ref="favorites-btn"
                >
                    <i class="bi bi-heart"></i>
                </button>
            </RouterLink>
        </div>

        <div class="content" v-show="!store.sidebarClosed">
            <div class="survey" v-if="survey.show.value">
                <p>Have any suggestions? Two questions, 100% anonymous:</p>
                <p>
                    <a :href="SURVEY_LINK" target="_blank">{{ SURVEY_LINK }}</a>
                </p>
                <p>
                    <i class="bi bi-info-circle"></i> The survey can also be
                    found in the "about" tab.
                </p>
                <p>
                    <button
                        class="btn-primary btn-rounded"
                        @click="survey.onHide"
                    >
                        don't show again
                    </button>
                </p>
            </div>
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
    padding: 0;
    width: $sidebar-width;
    height: $sidebar-width;
    color: $color-primary-lighter;

    .bi {
        font-size: 48px;
    }
}

.btn-small {
    width: $sidebar-width - 10px;

    &:hover {
        color: $color-primary-light;
    }

    .bi {
        font-size: 30px;
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
