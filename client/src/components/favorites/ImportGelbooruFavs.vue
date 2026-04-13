<script setup lang="ts">
import { computed, ref } from "vue";
import { delay } from "@/fetch";
import store from "@/store";
import ScreenCover from "../ScreenCover.vue";

defineEmits(["close"]);
const input = ref("");
const importing = ref(false);
const importPageMax = ref(0);
const importPage = ref(0);

const btnText = computed(() => {
    if (!importing.value) {
        return "Import";
    }

    return `Importing (${importPage.value} of ${importPageMax.value})`;
});

function onInput(e: Event) {
    const $el = e.target as HTMLInputElement;
    const value = $el.value;
    input.value = value.replace(/[^\d]+/, "");
    // Update the DOM in case input.value didn't change
    $el.value = input.value;
}

async function doImport() {
    importPage.value = 1;
    let totalResults = 0;
    const prevFavCount = store.favoritePosts().value.length;

    while (true) {
        const resp = await store.searchPosts(
            { include: [`fav:${input.value}`], exclude: [] },
            importPage.value,
        );
        totalResults = resp.total_count;
        importPageMax.value = Math.ceil(totalResults / resp.count_per_page);

        if (!resp.results.length) {
            break;
        }

        await store.addFavoritePosts(resp.results);

        // Last page
        if (resp.results.length < resp.count_per_page) {
            break;
        }

        // Avoid hitting rate limit
        await delay(1000);
        importPage.value++;
    }

    const added = store.favoritePosts().value.length - prevFavCount;
    let msg = "";

    if (totalResults === 0) {
        msg = "No favs found";
    } else if (added === 0) {
        msg = "No new favs to import";
    } else if (added === 1) {
        msg = "Imported 1 fav";
    } else {
        msg = `Imported ${added} favs`;
    }

    store.toast = {
        msg,
        type: "info",
    };
    // bonus: show progress modal
}

async function onSubmit() {
    if (importing.value) {
        return;
    }

    importing.value = true;

    try {
        await doImport();
    } finally {
        importing.value = false;
    }
}
</script>

<template>
    <div class="modal-container">
        <ScreenCover @click="$emit('close')" />

        <div class="modal">
            <p>
                To import your favs from Gelbooru, you'll need your account ID.
            </p>
            <ol>
                <li>Login to your Gelbooru account</li>
                <li>
                    Go to the
                    <a
                        href="Import your favorited posts from Gelbooru."
                        target="_blank"
                        >account page</a
                    >
                </li>
                <li>Click on <strong>My Profile</strong></li>
                <li>Copy the number in the URL</li>
            </ol>
            <form @submit.prevent="onSubmit">
                <input
                    type="text"
                    class="text-input rounded"
                    :value="input"
                    @input="onInput"
                    placeholder="Gelbooru ID"
                />
                <button
                    type="submit"
                    class="btn-primary btn-rounded"
                    :disabled="importing"
                >
                    {{ btnText }}
                </button>
            </form>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.modal-container {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1;
}

.modal {
    position: relative;
    z-index: 1000;
    color: white;
    padding: 0.8em;
    background-color: $color-gray;
}

.text-input {
    background-color: $color-darkgray;
}
</style>
