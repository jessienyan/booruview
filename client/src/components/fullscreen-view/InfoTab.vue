<script setup lang="ts">
import store from "@/store";
import { computed, ref, watchEffect } from "vue";
import TagList from "../TagList.vue";

const tags = ref<Tag[]>([]);
const { post } = defineProps<{ post: Post }>();

watchEffect(() => {
    if (post === null) {
        return;
    }

    store.tagsForPost(post).then((val) => (tags.value = val));
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
        }

        return ret;
    }),
);
</script>

<template>
    <div class="tag-list">
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
