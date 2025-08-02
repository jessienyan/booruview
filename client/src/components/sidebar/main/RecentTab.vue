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

function onDelete(index: number) {
    store.settings.queryHistory.splice(index, 1);
    store.saveSettings();
}

function onSearch(query: SearchQuery) {
    if (store.fetchingPosts) {
        return;
    }

    store.query = query.copy();
    store.currentPage = 1;
    store.clearPosts();

    store.searchPosts().catch(() => store.clearPosts());
}
</script>

<template>
    <div class="recent-list" ref="list">
        <div
            class="history-entry"
            v-for="(entry, i) of store.settings.queryHistory"
            :key="entry.date.getTime()"
        >
            <div class="tag-list">
                <TagList :tags="styledTags(entry.query).value" />
            </div>
            <div class="entry-footer">
                <button
                    class="btn-delete"
                    title="remove from history"
                    @click="onDelete(i)"
                >
                    <i class="bi bi-trash3"></i>
                </button>
                <span class="time" :title="entry.date.toLocaleString()">{{
                    relativeTime(entry.date)
                }}</span>
                <button
                    class="btn-primary btn-rounded btn-search"
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
    margin-bottom: 1rem;
}

.entry-footer {
    border-top: 1px solid $color-primary-lighter;
    background-color: $color-primary-darker;
    display: flex;
    flex-direction: row;
    align-items: center;

    .time {
        font-size: 16px;
        color: $color-primary-light;
    }
}

.btn-search {
    margin: 0.8rem 0.8rem 0.8rem auto;
}

.tag-list {
    padding: 0.8rem;
    background-color: $color-primary-darker;
}

.btn-delete {
    border: none;
    background: none;
    cursor: pointer;
    padding: 0.6rem;
    margin: 0.6rem;
    font-size: 20px;
    color: $color-primary-light;
}
</style>
