<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import Footer from "@/components/Footer.vue";
import PageChangeGesture from "@/PageChangeGesture.vue";
import PostContainer from "@/components/PostContainer.vue";
import { useMainContainer } from "@/composable";
import { onMounted } from "vue";

const mainContainer = useMainContainer();

onMounted(() => mainContainer.value.focus());
</script>

<template>
    <PageChangeGesture :scroll-container="mainContainer" />
    <NoResults v-if="store.totalPostCount === 0 && store.hasSearched">
        no results :(
    </NoResults>
    <template v-else-if="store.totalPostCount > 0">
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="mainContainer"
        />
        <Footer />
    </template>
</template>
