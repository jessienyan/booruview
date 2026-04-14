<script setup lang="ts">
import { onMounted } from "vue";
import Footer from "@/components/Footer.vue";
import LoadingResults from "@/components/LoadingResults.vue";
import NoResults from "@/components/NoResults.vue";
import PostContainer from "@/components/PostContainer.vue";
import { useKeydownListener, useMainContainer } from "@/composable";
import PageSwipeArrow from "@/PageSwipeArrow.vue";
import store from "@/store";

const mainContainer = useMainContainer();
onMounted(() => mainContainer.value.focus());
useKeydownListener("ArrowLeft", mainContainer, () => store.prevPage());
useKeydownListener("ArrowRight", mainContainer, () => store.nextPage());
</script>

<template>
    <LoadingResults v-if="store.fetchingPosts" />
    <PageSwipeArrow
        :scroll-container="mainContainer"
        :current-page="store.currentPage"
        :max-page="store.maxPage()"
        @prev="store.prevPage()"
        @next="store.nextPage()"
    />
    <NoResults v-if="store.totalPostCount === 0 && store.hasSearched">
        no results :(
    </NoResults>
    <template v-else-if="store.totalPostCount > 0">
        <!-- Using keyed=false to save time on garbage collection when changing pages -->
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="mainContainer"
            :keyed="false"
        />
        <p
            v-if="store.maxPage() > 200 && store.currentPage >= 200"
            class="end-notice"
        >
            Unfortunately, results past page 200 aren't viewable<br />because
            they are blocked by Gelbooru. :(
        </p>
        <Footer
            v-else
            :current-page="store.currentPage"
            :max-page="store.maxPage()"
            :total-count="store.totalPostCount"
            :prev-to="{
                name: 'search',
                params: {
                    page: (store.currentPage - 1).toString(),
                    query: $route.params.query,
                },
            }"
            :next-to="{
                name: 'search',
                params: {
                    page: (store.currentPage + 1).toString(),
                    query: $route.params.query,
                },
            }"
            :prev-disabled="store.fetchingPosts"
            :next-disabled="store.fetchingPosts || store.currentPage >= 200"
        />
    </template>
</template>
