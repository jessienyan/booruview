<script setup lang="ts">
import store from "@/store";
import Sidebar from "@/components/sidebar/Sidebar.vue";
import { RouterView } from "vue-router";
</script>

<template>
    <div
        class="app"
        :class="{
            'sidebar-closed': store.sidebarClosed,
            'sidebar-open': !store.sidebarClosed,
        }"
    >
        <Sidebar />
        <main
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
