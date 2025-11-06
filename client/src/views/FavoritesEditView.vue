<script setup lang="ts">
import store from "@/store";
import DraggablePostContainer from "@/components/DraggablePostContainer.vue"
import { useMainContainer } from "@/composable";

const mainContainer = useMainContainer();
</script>

<template>
    <DraggablePostContainer :posts="store.settings.favorites" :scroll-container="mainContainer" />
</template>
