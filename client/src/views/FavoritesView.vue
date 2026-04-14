<script setup lang="ts">
import { onMounted } from "vue";
import FavoritesHeader from "@/components/favorites/FavoritesHeader.vue";
import NoResults from "@/components/NoResults.vue";
import PostContainer from "@/components/PostContainer.vue";
import { useMainContainer } from "@/composable";
import store from "@/store";

const favPosts = store.favoritePosts();
const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());
</script>

<template>
    <FavoritesHeader />
    <NoResults v-if="favPosts.length === 0">
        you don't have any favorites yet
    </NoResults>
    <template v-else>
        <div class="spacer"></div>
        <PostContainer
            :posts="favPosts"
            :scroll-container="mainContainer"
            :keyed="true"
        />
    </template>
</template>

<style lang="css" scoped>
.spacer {
    height: 1rem;
}
</style>
