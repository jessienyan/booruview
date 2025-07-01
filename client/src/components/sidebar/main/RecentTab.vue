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
</script>

<template>
    <div class="recent-list">
        <div class="history-entry" v-for="entry of store.settings.queryHistory">
            <p>{{ relativeTime(entry.date) }}</p>
            <TagList :tags="styledTags(entry.query).value" />
        </div>
    </div>
</template>
