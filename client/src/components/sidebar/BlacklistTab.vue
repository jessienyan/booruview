<script setup lang="ts">
import TagList from "../TagList.vue";
import store from "@/store";
</script>

<template>
    <TagList
        v-if="store.settings.blacklist.length > 0"
        :jiggle="false"
        :show-checkmark="false"
        :tags="store.settings.blacklist"
    />
    <template v-else>
        <p>
            Any tags you've blacklisted will appear here. Posts that have
            blacklisted tags will be hidden.
        </p>
        <p>To blacklist, click a tag and then select "blacklist".</p>
    </template>
</template>
