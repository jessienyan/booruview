<script setup lang="ts">
import store from "@/store";
import { computed, ref, watchEffect, type CSSProperties } from "vue";
import TagList from "../TagList.vue";

const tags = ref<Tag[]>([]);
const { post } = defineProps<{ post: Post }>();

watchEffect(() => {
    if (post === null) {
        return;
    }

    store.tagsForPost(post).then((val) => (tags.value = val));
});

// Add padding if the menu would cover part of the container
const containerStyle = computed<CSSProperties>(() => {
    if (!store.settings.fullscreenViewMenuRotate) {
        return {};
    }

    switch (store.settings.fullscreenViewMenuAnchor) {
        case "bottomleft":
        case "left":
        case "topleft":
            return {
                paddingLeft: "3.5rem",
            };

        case "bottomright":
        case "right":
        case "topright":
            return {
                paddingRight: "3.5rem",
            };
    }

    return {};
});

const styledTags = computed(() =>
    tags.value.map((t) => {
        const ret: TagChip = {
            tag: t,
            style: "default",
        };

        if (store.query.isIncluded(t.name)) {
            ret.style = "checkmark";
        } else if (store.query.isExcluded(t.name)) {
            ret.style = "strikethrough";
        } else if (
            store.settings.blacklist.findIndex((bl) => bl.name === t.name) !==
            -1
        ) {
            ret.style = "blacklist";
        }

        return ret;
    }),
);
</script>

<template>
    <div class="tag-list" :style="containerStyle">
        <TagList :jiggle="false" :tags="styledTags" />
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/breakpoints";
@import "@/assets/mixin";

.tag-list {
    max-width: 800px;
    overflow-y: scroll;
    max-height: 100%;

    @include hide-scrollbar;
}
</style>
