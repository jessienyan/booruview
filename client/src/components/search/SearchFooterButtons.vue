<script setup lang="ts">
import { computed } from "vue";
import store from "@/store";

const savedSearches = store.savedSearches();
const isSaved = computed(() => {
    const index = savedSearches.value.findIndex((saved) =>
        saved.query.equals(store.query),
    );
    return index !== -1;
});

function clearTags() {
    store.query.clear();
}

function saveSearch() {
    store.addToSavedSearches([
        {
            date: new Date(),
            query: store.query.copy(),
        },
    ]);
}

function unsaveSearch() {
    store.removeFromSavedSearches([store.query.toJSONSimple()]);
}
</script>

<template>
    <div class="footer-btns">
        <button title="clear tags" @click="clearTags()">
            <span class="bi bi-eraser-fill"></span>
        </button>
        <button v-if="!isSaved" title="save search" @click="saveSearch()">
            <span class="bi bi-bookmark"></span>
        </button>
        <button v-if="isSaved" title="unsave search" @click="unsaveSearch()">
            <span class="bi bi-bookmark-fill"></span>
        </button>
    </div>
</template>

<style lang="scss" scoped>
.footer-btns {
    display: flex;
    justify-content: flex-end;
    gap: 0.8em;
    margin-top: 0.8em;

    button {
        background-color: #333;
        border: 1px solid #555;
        color: #ddd;
        border-radius: 4px;
        padding: 0.3em 0.4em;
        width: 36px;
        height: 36px;
        font-size: 1.2em;
        cursor: pointer;
    }
}
</style>
