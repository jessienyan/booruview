<script setup lang="ts">
import store from "@/store";
import { ref, computed, useTemplateRef } from "vue";
import DropdownMenu from "../DropdownMenu.vue";
import Chip from "../tag-chip/Chip.vue";

const showFavorites = ref(false);
const hasFavorites = computed(() => store.settings.favoriteTags.length > 0);
const favoritesBtn = useTemplateRef("favoritesBtn");

const categoryOrder: TagType[] = [
    "artist",
    "character",
    "copyright",
    "tag",
    "unknown",
    "metadata",
    "deprecated",
];

// Favorited tags sorted by category then by name
const tags = computed<TagChip[]>(() => {
    const sorted = [...store.settings.favoriteTags].sort((a, b) => {
        const category =
            categoryOrder.indexOf(a.type) - categoryOrder.indexOf(b.type);
        if (category !== 0) {
            return category;
        }
        return a.name.localeCompare(b.name);
    });
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
        v-if="hasFavorites"
        ref="favoritesBtn"
        class="btn-favorite-tags btn-primary"
        @click="showFavorites = !showFavorites"
    >
        <i class="bi bi-heart-fill"></i>{{ " "
        }}<i
            class="bi"
            :class="{
                'bi-caret-down-fill': !showFavorites,
                'bi-caret-up-fill': showFavorites,
            }"
        ></i>
    </button>
    <DropdownMenu :el="favoritesBtn" v-model:show="showFavorites">
        <div class="fav-tags">
            <Chip
                v-for="t of tags"
                :tag="t"
                :can-edit="false"
                :show-heart="false"
            />
        </div>
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-favorite-tags {
    border-top-right-radius: 4px;
    border-bottom-right-radius: 4px;
}

.fav-tags {
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
