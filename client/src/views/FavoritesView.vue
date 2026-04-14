<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import Footer from "@/components/Footer.vue";
import NoResults from "@/components/NoResults.vue";
import PostContainer from "@/components/PostContainer.vue";
import { useMainContainer } from "@/composable";
import { POSTS_PER_PAGE } from "@/config";
import store from "@/store";

const route = useRoute();
const favPosts = store.favoritePosts();
const mainContainer = useMainContainer();
const currentPage = computed(() =>
    parseInt((route.params.page as string) || "1", 10),
);
const maxPage = computed(() =>
    Math.ceil(favPosts.value.length / POSTS_PER_PAGE),
);
const currentPosts = computed(() =>
    favPosts.value.slice(
        (currentPage.value - 1) * POSTS_PER_PAGE,
        currentPage.value * POSTS_PER_PAGE,
    ),
);

onMounted(() => mainContainer.value.focus());
</script>

<template>
    <NoResults v-if="favPosts.length === 0">
        you don't have any favorites yet
    </NoResults>
    <template v-else>
        <PostContainer
            :posts="currentPosts"
            :scroll-container="mainContainer"
            :keyed="false"
        />
        <Footer
            :current-page="currentPage"
            :max-page="maxPage"
            :total-count="favPosts.length"
            :prev-to="{
                name: 'favorites',
                params: {
                    page: (currentPage - 1).toString(),
                },
            }"
            :next-to="{
                name: 'favorites',
                params: {
                    page: (currentPage + 1).toString(),
                },
            }"
            :prev-disabled="currentPage === 1"
            :next-disabled="currentPage >= maxPage"
        />
    </template>
</template>
