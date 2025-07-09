<script setup lang="ts">
import SidebarFooter from "./footer/SidebarFooter.vue";
import SidebarMain from "./main/SidebarMain.vue";

defineEmits(["toggle"]);
const { closed } = defineProps<{ closed: boolean }>();
</script>

<template>
    <header class="sidebar-container">
        <div class="sidebar-header">
            <button class="toggle-btn" @click="$emit('toggle')">
                <i class="bi bi-list"></i>
            </button>
        </div>

        <div class="sidebar-content" v-show="!closed">
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

.toggle-btn {
    background: none;
    border: none;
    font-size: 40px;
    cursor: pointer;

    color: $color-primary-lighter;
}
</style>
