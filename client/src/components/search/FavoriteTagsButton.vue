<script setup lang="ts">
import store from "@/store";
import { computed } from "vue";
import CollapsableTagMenu from "@/components/CollapsableTagMenu.vue";

const tags = computed<TagChip[]>(() =>
    store.settings.favoriteTags.map<TagChip>(tag => {
        let style: ChipStyle = "default";

        if (store.query.isIncluded(tag.name)) {
            style = "checkmark";
        } else if (store.query.isExcluded(tag.name)) {
            style = "strikethrough";
        }

        return {
            tag,
            style
        };
    }));
</script>

<template>
    <CollapsableTagMenu :tags="tags"><i class="bi bi-heart-fill"></i>{{ " " }}</CollapsableTagMenu>
</template>
