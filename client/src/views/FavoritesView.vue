<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import DraggablePostContainer from "@/components/DraggablePostContainer.vue";
import { onMounted } from "vue";
import { useDontShowAgain, useMainContainer } from "@/composable";

const desktopNotice = useDontShowAgain("desktop-fav-experimental-notice");

if(desktopNotice.show.value) {
    store.toast = {
        msg: "new: drag-and-drop favorites to sort them. DESKTOP ONLY & experimental. mobile support coming soon",
        type: "info",
    }
    desktopNotice.ack();
}

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
