<script setup lang="ts">
import store from "@/store";
import NoResults from "@/components/NoResults.vue";
import Footer from "@/components/Footer.vue";
import PostContainer from "@/components/PostContainer.vue";
import { onMounted, watch } from "vue";
import { tagsToSearchQuery } from "@/search";

const props = defineProps<{
    page: number | string;
    query?: string | string[];
}>();

function loadPosts() {
    let page: number;

    if (typeof props.page === "string") {
        page = parseInt(props.page);
    } else {
        page = props.page;
    }

    tagsToSearchQuery(props.query || []).then((q) => {
        store.query = q;
        store.searchPosts(page);
    });
}

watch(() => [props.page, props.query], loadPosts);
onMounted(loadPosts);
</script>

<template>
    <NoResults v-if="store.totalPostCount === 0"> no results :( </NoResults>
    <template v-else>
        <PostContainer
            :posts="store.postsForCurrentPage() || []"
            :scroll-container="$root?.$parent?.$el"
        />
        <Footer />
    </template>
</template>
