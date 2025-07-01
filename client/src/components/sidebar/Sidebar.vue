<script setup lang="ts">
import { useDontShowAgain } from "@/composable";
import SidebarFooter from "./footer/SidebarFooter.vue";
import SidebarMain from "./main/SidebarMain.vue";

defineEmits(["toggle"]);
const { closed } = defineProps<{ closed: boolean }>();
const survey = useDontShowAgain("hide-survey");
</script>

<template>
    <header class="sidebar-container">
        <div class="sidebar-header">
            <button class="toggle-btn" @click="$emit('toggle')">
                <i class="bi bi-list"></i>
            </button>
        </div>

        <div class="sidebar-content" v-show="!closed">
            <div class="survey" v-if="survey.show.value">
                <p>
                    Do you like booruview or hate it? Let me know with this 3
                    question survey. Anonymous, no signup.
                </p>
                <p>
                    <a
                        href="https://freesurveys.org/s/afOcVj9z-f"
                        target="_blank"
                        >https://freesurveys.org/s/afOcVj9z-f</a
                    >
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

.sidebar-content {
    width: 450px;
    margin-top: 10px;
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

.toggle-btn {
    background: none;
    border: none;
    font-size: 40px;
    cursor: pointer;

    color: $color-primary-lighter;
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
