<script setup lang="ts">
import { type CSSProperties, computed, ref, watchEffect } from "vue";
import store from "@/store";
import TagList from "../TagList.vue";

const tags = ref<Tag[]>([]);
const { post } = defineProps<{ post: Post }>();

watchEffect(() => {
    if (post === null) {
        return;
    }

    store.tagsForPost(post).then((val) => {
        tags.value = val;
    });
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

const blacklist = store.blacklist();

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
            blacklist.value.findIndex((bl) => bl.name === t.name) !== -1
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
        <p>
            <a :href="post.image_url || post.lowres_url" target="_blank"
                >download<i class="bi bi-download"></i
            ></a>
            -
            <a
                :href="`https://gelbooru.com/index.php?page=post&s=view&id=${post.id}`"
                target="_blank"
                >view on Gelbooru<i
                    class="external-link bi bi-box-arrow-up-right"
                ></i
            ></a>
        </p>
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

p {
    margin-top: 30px;
    text-align: right;
}

a {
    font-weight: bold;

    .bi {
        margin-left: 6px;
        vertical-align: baseline;
    }
}
</style>
