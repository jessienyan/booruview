<script setup lang="ts">
import { computed, useTemplateRef } from "vue";
import store from "@/store";
import DropdownMenu from "../DropdownMenu.vue";

const btnRef = useTemplateRef("button");
const showMenu = defineModel<boolean>({ default: false });

const includeCount = computed(() => store.query._include.size);
const excludeCount = computed(() => store.query._exclude.size);

function clearIncluded() {
    store.query._include.clear();
}

function clearExcluded() {
    store.query._exclude.clear();
}
</script>

<template>
    <button
        ref="button"
        class="btn-gray btn-rounded btn-clear-tags"
        @click="showMenu = !showMenu"
    >
        clear tags
        <i
            class="bi"
            :class="{
                'bi-caret-down-fill': !showMenu,
                'bi-caret-up-fill': showMenu,
            }"
        ></i>
    </button>

    <DropdownMenu :el="btnRef" v-model:show="showMenu">
        <button
            v-if="includeCount > 0"
            class="dropdown-option btn-gray"
            @click="clearIncluded"
        >
            included ({{ includeCount }})
        </button>
        <button
            v-if="excludeCount > 0"
            class="dropdown-option btn-gray"
            @click="clearExcluded"
        >
            excluded ({{ excludeCount }})
        </button>
        <button
            class="dropdown-option btn-gray"
            @click="
                clearIncluded();
                clearExcluded();
            "
        >
            all
        </button>
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-clear-tags {
    margin-top: 0.8rem;
    margin-left: auto;
    display: block;
}
</style>
