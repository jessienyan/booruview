<script setup lang="ts">
import { ref, Teleport } from "vue";
import store from "@/store";
import ImportGelbooruFavs from "./ImportGelbooruFavs.vue";

const showFavsModal = ref(false);

function onToggle() {
    store.settings.favHeaderOpen = !store.settings.favHeaderOpen;
    store.saveSettings();
}
</script>

<template>
    <div class="header-container" v-if="store.settings.favHeaderOpen">
        <button class="btn-primary btn-rounded" @click="showFavsModal = true">
            import gelbooru favs
        </button>
    </div>
    <button
        class="header-toggle-btn"
        :class="{ closed: !store.settings.favHeaderOpen }"
        @click="onToggle"
    >
        <i v-if="store.settings.favHeaderOpen" class="bi bi-chevron-up"></i>
        <i v-else class="bi bi-chevron-down"></i>
    </button>

    <Teleport to="body">
        <ImportGelbooruFavs
            v-if="showFavsModal"
            @close="showFavsModal = false"
        />
    </Teleport>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";

.header-container {
    position: relative;
    background-color: $color-primary-darker;
    padding: 0.8em;
    box-shadow: 0 0 10px $color-bg;
}

.header-toggle-btn {
    border: 1px solid #444;
    border-top: none;
    padding: 0.4rem 1.6rem;
    border-radius: 0 0 4px 4px;
    background-color: #1e1e1e;
    color: #888;
    cursor: pointer;
    display: block;
    margin: auto;
    width: fit-content;

    &.closed {
        opacity: 0.6;
    }
}
</style>
