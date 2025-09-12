<script setup lang="ts">
import store from "@/store";
import { ref, computed, useTemplateRef } from "vue";
import DropdownMenu from "../DropdownMenu.vue";

const showFavorites = ref(false);
const hasFavorites = computed(() => store.settings.favoriteTags.length > 0);
const favoritesBtn = useTemplateRef("favoritesBtn");
</script>

<template>
    <button
        v-if="hasFavorites"
        class="btn-favorite-tags btn-primary"
        ref="favoritesBtn"
    >
        <i class="bi bi-heart-fill"></i>{{ " "
        }}<i
            class="bi"
            :class="{
                'bi-caret-down-fill': true,
                'bi-caret-up-fill': false,
            }"
        ></i>
    </button>
    <DropdownMenu :el="favoritesBtn" v-model:show="showFavorites">
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-favorite-tags {
    border-left: none;
}
</style>
