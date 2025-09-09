<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import PostContainer from "@/components/PostContainer.vue";
</script>

<template>
    <NoResults v-if="store.settings.favorites.length === 0">
        you don't have any favorites yet
    </NoResults>
    <PostContainer
        v-else
        :posts="store.settings.favorites"
        :scroll-container="$root?.$parent?.$el"
    />
</template>
