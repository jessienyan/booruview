<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import PostContainer from "@/components/PostContainer.vue";
import { onMounted } from "vue";
import { useMainContainer } from "@/composable";

const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());
</script>

<template>
    <NoResults v-if="store.settings.favorites.length === 0">
        you don't have any favorites yet
    </NoResults>
    <PostContainer
        v-else
        :posts="store.settings.favorites"
        :scroll-container="mainContainer"
    />
</template>
