<script setup lang="ts">
import { computed, onUnmounted, ref } from "vue";
import { delay } from "@/fetch";
import store from "@/store";
import ScreenCover from "../ScreenCover.vue";

const emit = defineEmits(["close"]);
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

    while (importing.value) {
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
        msg = "No favorites found (is the account ID right?)";
    } else if (added === 0) {
        msg = "No new favorites to import";
    } else if (added === 1) {
        msg = "Done! Imported 1 favorite";
    } else {
        msg = `Done! Imported ${added} favorites`;
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

function onClose() {
    importing.value = false;
    emit("close");
}

onUnmounted(() => {
    importing.value = false;
});
</script>

<template>
    <div class="modal-container">
        <ScreenCover @click="onClose" />

        <div class="modal">
            <p>
                To import your favorite posts from Gelbooru, you'll need your
                account ID.
            </p>
            <ol>
                <li>Login to your Gelbooru account</li>
                <li>
                    Go to the
                    <a
                        href="https://gelbooru.com/index.php?page=account&s=home"
                        target="_blank"
                        >account page</a
                    >
                </li>
                <li>Click on <strong>My Profile</strong></li>
                <li>Copy the number in the URL, for example:</li>
            </ol>
            <pre>https://gelbooru.com/...&amp;id=<strong>123456789</strong></pre>
            <form @submit.prevent="onSubmit">
                <input
                    type="text"
                    class="text-input rounded input-block"
                    :value="input"
                    @input="onInput"
                    placeholder="account id here"
                />
                <button
                    type="submit"
                    class="btn-primary btn-rounded btn-block"
                    :disabled="importing || !input.length"
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
    color: #ccc;
    padding: 0.01em 0.8em;
    margin: 1rem;
    max-width: 400px;
    background-color: $color-primary;
}

.text-input {
    background-color: $color-darkgray;
}

form {
    border-top: 2px solid $color-primary-lighter;
    margin-top: 1em;
}

pre {
    background-color: $color-darkgray;
    color: #999;
    padding: 1em;
    font-size: 0.8em;
    width: min-content;
    margin: auto;
    pointer-events: none;
    user-select: none;

    strong {
        color: white;
    }
}

ol {
    padding-left: 2em;
}

li {
    margin-bottom: 0.5em;
}
</style>
