<script setup lang="ts">
import { computed } from "vue";
import TagList from "@/components/TagList.vue";
import { useRelativeTime } from "@/composable";
import type { SearchQuery } from "@/search";
import store from "@/store";

const relativeTime = useRelativeTime();
const history = store.searchHistory();

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
    store.removeFromSearchHistory(history.value[index]);
}
</script>

<template>
    <div class="recent-list" ref="list">
        <div
            class="history-entry"
            v-for="(entry, i) of history"
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
                <RouterLink
                    class="btn-search"
                    :to="{
                        name: 'search',
                        params: { page: 1, query: entry.query.asQueryParams() },
                        force: true,
                    }"
                >
                    <button class="btn-primary btn-rounded">search</button>
                </RouterLink>
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";

.history-entry {
    margin-bottom: 0.8rem;
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
