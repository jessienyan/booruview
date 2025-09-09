<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import Footer from "@/components/Footer.vue";
import PageChangeGesture from "@/PageChangeGesture.vue";
import PostContainer from "@/components/PostContainer.vue";
import { onMounted } from "vue";
import { tagsToSearchQuery } from "@/search";
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRoute } from "vue-router";
import { useMainContainer } from "@/composable";

const mainContainer = useMainContainer();
const route = useRoute();

function loadPosts(page: number, query: string | string[]) {
    return new Promise<void>((resolve, reject) => {
        tagsToSearchQuery(query || []).then((q) => {
            store.query = q;
            store
                .searchPosts(page)
                .then(() => {
                    mainContainer.value.focus();
                    resolve();
                })
                .catch(reject);
        });
    });
}

onBeforeRouteUpdate((to) => {
    const { page, query } = to.params;
    return loadPosts(parseInt(page as string), query);
});

onBeforeRouteLeave((to, from) => {
    store.lastSearchRoute = from;
});

onMounted(() => {
    const { page, query } = route.params;
    store.lastSearchRoute = route;
    loadPosts(parseInt(page as string), query);
});
</script>

<template>
    <PageChangeGesture :scroll-container="mainContainer" />
    <NoResults v-if="store.totalPostCount === 0 && !store.fetchingPosts">
        no results :(
    </NoResults>
    <template v-else-if="!store.fetchingPosts">
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="mainContainer"
        />
        <Footer />
    </template>
</template>
