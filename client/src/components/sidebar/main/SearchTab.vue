<script lang="ts" setup>
import { computed } from "vue";
import { useRouter } from "vue-router";
import ClearTagsButton from "@/components/search/ClearTagsButton.vue";
import SearchForm from "@/components/search/SearchForm.vue";
import TagList from "@/components/TagList.vue";
import store from "@/store";

const router = useRouter();

function doPostSearch() {
    if (store.fetchingPosts) {
        return;
    }

    const searchDidntChange =
        router.currentRoute.value.name === "search" &&
        store.currentPage === 1 &&
        store.lastQuery.equals(store.query);

    store.justClickedSearchButton = true;

    router.push({
        name: "search",
        params: { page: 1, query: store.query.asQueryParams() },
        force: true,
        // Searching while nothing has changed is the same as refreshing; don't store a new entry
        replace: searchDidntChange,
    });
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

        <div v-if="styledTags.length" class="taglist-container">
            <TagList
                :jiggle="true"
                :tags="styledTags"
                :actions="{ edit: true }"
            />
            <ClearTagsButton />
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

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

.clear-tags {
    margin-left: auto;
    display: block;
}
</style>
