<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import Footer from "@/components/Footer.vue";
import PageSwipeArrow from "@/PageSwipeArrow.vue";
import PostContainer from "@/components/PostContainer.vue";
import { useMainContainer } from "@/composable";
import { onMounted } from "vue";
import LoadingResults from "@/components/LoadingResults.vue";

const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());
</script>

<template>
    <LoadingResults v-if="store.fetchingPosts" />
    <PageSwipeArrow :scroll-container="mainContainer" />
    <NoResults v-if="store.totalPostCount === 0 && store.hasSearched">
        no results :(
    </NoResults>
    <template v-else-if="store.totalPostCount > 0">
        <!-- Using keyed=false to save time on garbage collection when changing pages -->
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="mainContainer"
            :keyed=false
        />
        <Footer />
    </template>
</template>
