<script setup lang="ts">
import store from "@/store";
import { ref, computed, useTemplateRef } from "vue";
import DropdownMenu from "../DropdownMenu.vue";
import Chip from "../tag-chip/Chip.vue";
import { sortTags } from "@/tag";

const open = ref(false);
const hasFavorites = computed(() => store.settings.favoriteTags.length > 0);
const btnRef = useTemplateRef("button");

// Favorited tags sorted by category then by name
const tags = computed<TagChip[]>(() => {
    const sorted = sortTags(store.settings.favoriteTags);
    const styled: TagChip[] = [];

    for (const tag of sorted) {
        let style: ChipStyle = "default";

        if (store.query.isIncluded(tag.name)) {
            style = "checkmark";
        } else if (store.query.isExcluded(tag.name)) {
            style = "strikethrough";
        }

        styled.push({
            tag,
            style,
        });
    }

    return styled;
});
</script>

<template>
    <button
        ref="button"
        class="btn-quick-tags btn-primary"
        @click="open = !open"
    >
        <i class="bi bi-tags-fill"></i>{{ " "
        }}<i
            class="bi"
            :class="{
                'bi-caret-down-fill': !open,
                'bi-caret-up-fill': open,
            }"
        ></i>
    </button>
    <DropdownMenu :el="btnRef" v-model:show="open">
        <div class="chip-list">
            <Chip v-for="t of tags" :tag="t" :show-heart="false" />
        </div>
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-quick-tags {
    border-top-right-radius: 4px;
    border-bottom-right-radius: 4px;
}

.chip-list {
    display: flex;
    flex-direction: column;
    background-color: #1c1c1c;
    gap: 8px;
    padding: 8px;
    max-width: 300px;
    max-height: 350px;
    overflow-y: scroll;

    &:deep(.chip) {
        margin: 0;
        padding: 0.8rem;
    }
}
</style>
