<script setup lang="ts">
import { computed } from "vue";
import TagList from "@/components/TagList.vue";
import store from "@/store";

const styledTags = computed<TagChip[]>(() =>
    store.settings.blacklist.map((tag) => ({ tag, style: "default" })),
);
</script>

<template>
    <TagList
        v-if="store.settings.blacklist.length > 0"
        :jiggle="false"
        :tags="styledTags"
    />
    <template v-else>
        <p>
            Any tags you've blacklisted will appear here. Posts that have
            blacklisted tags will be hidden.
        </p>
        <p>To blacklist, click a tag and then select "blacklist".</p>
    </template>
</template>
