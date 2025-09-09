<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import Footer from "@/components/Footer.vue";
import PageChangeGesture from "@/PageChangeGesture.vue";
import PostContainer from "@/components/PostContainer.vue";
import { onMounted, watch } from "vue";
import { tagsToSearchQuery } from "@/search";

const props = defineProps<{
    page: string;
    query?: string | string[];
}>();

function loadPosts() {
    tagsToSearchQuery(props.query || []).then((q) => {
        store.query = q;
        store.searchPosts(parseInt(props.page));
    });
}

watch(() => [props.page, props.query], loadPosts);
onMounted(loadPosts);
</script>

<template>
    <PageChangeGesture :scroll-container="$root?.$parent?.$el" />
    <NoResults v-if="store.totalPostCount === 0 && !store.fetchingPosts">
        no results :(
    </NoResults>
    <template v-else-if="!store.fetchingPosts">
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="$root?.$parent?.$el"
        />
        <Footer />
    </template>
</template>
