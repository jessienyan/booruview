<script setup lang="ts">
import { onMounted } from "vue";
import DraggablePostContainer from "@/components/DraggablePostContainer.vue";
import NoResults from "@/components/NoResults.vue";
import { useMainContainer } from "@/composable";
import store from "@/store";

const favPosts = store.favoritePosts();
const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());
</script>

<template>
    <NoResults v-if="favPosts.length === 0">
        you don't have any favorites yet
    </NoResults>
    <DraggablePostContainer
        v-else
        :posts="favPosts"
        :scroll-container="mainContainer"
        @change="(posts) => store.setFavoritePosts(posts)"
    />
</template>
