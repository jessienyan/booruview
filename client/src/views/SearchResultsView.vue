<script setup lang="ts">
import { onMounted } from "vue";
import Footer from "@/components/Footer.vue";
import LoadingResults from "@/components/LoadingResults.vue";
import NoResults from "@/components/NoResults.vue";
import PostContainer from "@/components/PostContainer.vue";
import { useMainContainer } from "@/composable";
import PageSwipeArrow from "@/PageSwipeArrow.vue";
import store from "@/store";

const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());
</script>
<template>
    <NoResults>
        Images aren't loading at the moment, working on a fix
    </NoResults>
    <!--
    <LoadingResults v-if="store.fetchingPosts" />
    <PageSwipeArrow :scroll-container="mainContainer" />
    <NoResults v-if="store.totalPostCount === 0 && store.hasSearched">
        no results :(
    </NoResults>
    <template v-else-if="store.totalPostCount > 0">
        Using keyed=false to save time on garbage collection when changing pages
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="mainContainer"
            :keyed="false"
        />
        <Footer />
    </template>
    -->
</template>
