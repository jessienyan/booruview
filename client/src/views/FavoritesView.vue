<script setup lang="ts">
import { onMounted } from "vue";
import DraggablePostContainer from "@/components/DraggablePostContainer.vue";
import NoResults from "@/components/NoResults.vue";
import { useMainContainer } from "@/composable";
import store from "@/store";

const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());

function onChangeFavOrder(posts: Post[]) {
    store.settings.favorites = posts;
    store.saveSettings();
}
</script>

<template>
    <NoResults v-if="store.settings.favorites.length === 0">
        you don't have any favorites yet
    </NoResults>
    <DraggablePostContainer
        v-else
        :posts="store.settings.favorites"
        :scroll-container="mainContainer"
        @change="onChangeFavOrder"
    />
</template>
