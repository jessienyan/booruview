<script lang="ts" setup>
import { computed } from "vue";
import SearchForm from "@/components/search/SearchForm.vue";
import TagList from "@/components/TagList.vue";
import store from "@/store";

function doPostSearch() {
    if (store.fetchingPosts) {
        return;
    }

    store.currentPage = 1;
    store.clearPosts();

    store.searchPosts().catch(() => store.clearPosts());
}

function onTagSelect(tag: Tag, negated: boolean) {
    if (!negated) {
        store.query.includeTag(tag);
    } else {
        store.query.excludeTag(tag);
    }

    // Slight hack: try looking up unknown tags and replace it with the real version
    if (tag.type === "unknown") {
        store.loadTags([tag.name]).then(() => {
            const real = store.cachedTags.get(tag.name);

            if (real === undefined) {
                return;
            }

            store.query.removeTag(tag);

            if (!negated) {
                store.query.includeTag(real);
            } else {
                store.query.excludeTag(real);
            }
        });
    }
}

const styledTags = computed(() => {
    let ret: TagChip[];

    ret = store.query.includedList().map((tag) => ({ tag, style: "default" }));
    ret = ret.concat(
        store.query
            .excludedList()
            .map((tag) => ({ tag, style: "strikethrough" })),
    );

    return ret;
});
</script>

<template>
    <div class="search">
        <SearchForm
            @on-search="doPostSearch"
            @on-tag-select="onTagSelect"
            :show-spinner="store.fetchingPosts"
        />

        <div class="taglist-container">
            <TagList :jiggle="true" :tags="styledTags" />
        </div>
    </div>
</template>

<style lang="css" scoped>
.search {
    min-height: 0;
    display: flex;
    flex-direction: column;
}

.taglist-container {
    margin-top: 1rem;
    overflow-y: scroll;
    flex: 1;
}
</style>
