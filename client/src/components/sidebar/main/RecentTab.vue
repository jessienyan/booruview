<script setup lang="ts">
import TagList from "@/components/TagList.vue";
import { useRelativeTime } from "@/composable";
import type { SearchQuery } from "@/search";
import store from "@/store";
import { computed } from "vue";

const relativeTime = useRelativeTime();

function styledTags(query: SearchQuery) {
    return computed(() => {
        let ret: TagChip[];

        ret = query.includedList().map((tag) => ({ tag, style: "default" }));
        ret = ret.concat(
            query
                .excludedList()
                .map((tag) => ({ tag, style: "strikethrough" })),
        );

        return ret;
    });
}

function onSearch(query: SearchQuery) {
    store.query = query.copy();
    store.searchPosts();

    console.log("search", query.asList());
}
</script>

<template>
    <div class="recent-list">
        <div
            class="history-entry"
            v-for="entry of store.settings.queryHistory"
            :key="entry.date.getTime()"
        >
            <div class="tag-list">
                <TagList :tags="styledTags(entry.query).value" />
            </div>
            <div class="entry-footer">
                <span class="time" :title="entry.date.toLocaleString()">{{
                    relativeTime(entry.date)
                }}</span>
                <button
                    class="btn-primary btn-rounded"
                    @click="onSearch(entry.query)"
                >
                    search
                </button>
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";

.history-entry {
    margin-bottom: 10px;
}

.entry-footer {
    border-top: 1px solid $color-primary-lighter;
    background-color: $color-primary-darker;
    padding: 10px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;

    .time {
        font-size: 16px;
        color: $color-primary-light;
    }
}

.tag-list {
    padding: 10px;
    background-color: $color-primary-darker;
}
</style>
