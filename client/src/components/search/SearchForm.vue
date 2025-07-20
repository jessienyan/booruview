<script setup lang="ts">
import { computed, ref, useTemplateRef, watch } from "vue";
import SearchSuggestions from "./Suggestions.vue";
import store from "@/store";
import { useDismiss } from "@/composable";

type SearchResponse = {
    results: Tag[];
};
const debounceMs = 100;

const emit = defineEmits<{
    onSearch: [];
    onTagSelect: [value: Tag, negated: boolean];
}>();

const { showSpinner = false } = defineProps<{
    showSpinner?: boolean;
}>();

const containerRef = useTemplateRef("container");
const forceRenderKey = ref(0);
const inputVal = ref("");
const selectedIndex = ref(-1);
const suggestions = ref<Tag[]>([]);
const timer = ref();
const inputRef = useTemplateRef("input");
const showSuggestions = ref(false);

useDismiss([containerRef.value], () => (showSuggestions.value = false));

function doTagSearch(query: string) {
    if (smartSuggestionValue.value.length === 0) {
        return;
    }

    // Encoding the query prevents trailing whitespace from being stripped
    fetch("/api/tagsearch?q=" + encodeURIComponent(smartSuggestionValue.value))
        .then((resp) => {
            // Request took too long and results don't match the input, discard
            if (query !== inputVal.value) {
                return;
            }

            resp.json().then((json: SearchResponse) => {
                suggestions.value = json.results.filter(
                    // Don't suggest tags already added to the search
                    (t) =>
                        !store.query.isIncluded(t.name) &&
                        !store.query.isExcluded(t.name) &&
                        store.settings.blacklist.findIndex(
                            (bl) => t.name === bl.name,
                        ) === -1,
                );
            });
        })
        .catch((err) => console.error(err));
}

function changeSelection(direction: number) {
    direction += selectedIndex.value;

    if (direction < 0) {
        direction = suggestions.value.length - 1;
    } else if (direction >= suggestions.value.length) {
        direction = 0;
    }

    selectedIndex.value = direction;
}

// Computes the current word(s) under the cursor in the search. This allows the
// autocomplete to be more intelligent by only replacing the word(s) the user is entering.
// For example if the query was `{blue sky ~ sunse`, then "sunse" would be the current word
// boundary (assuming the cursor was at the end).
const wordBoundary = computed(() => {
    // Use the cursor position to find the word(s) that will be autocompleted
    const end = inputRef.value?.selectionEnd ?? 0;
    let start = end;

    if (inputRef.value === null) {
        return { start, end };
    }

    // Special operators that act as boundaries
    const OPERATORS = ["-", "{", "~ " /* note the space */];

    for (let i = start; i > 0; i--) {
        let match = false;

        for (const op of OPERATORS) {
            if (inputVal.value.slice(i - op.length, i) === op) {
                match = true;
                start = i;
                break;
            }
        }

        if (match) {
            break;
        }

        // If no match just decrement start
        start--;
    }

    return { start, end };
});

const smartSuggestionValue = computed(() => {
    const { start, end } = wordBoundary.value;
    return inputVal.value.slice(start, end);
});

function autoComplete() {
    if (suggestions.value.length === 0 || inputRef.value === null) {
        return;
    }

    const index = selectedIndex.value === -1 ? 0 : selectedIndex.value;
    let suggestionVal = suggestions.value[index].name;

    const bounds = wordBoundary.value;
    const before = inputVal.value.substring(0, bounds.start);
    const after = inputVal.value.substring(bounds.end);
    const spliced = before + suggestionVal + after;

    inputVal.value = spliced;
}

function onSuggestionClick(index: number) {
    selectedIndex.value = index;
    autoComplete();

    const isOR = /^-?\{/.test(inputVal.value);

    // Don't submit when building an OR query since the user may want to select more tags
    if (!isOR) {
        onSubmit();
    }

    inputRef.value?.focus();
}

function onInput(e: Event) {
    if (e.target === null) {
        return;
    }

    // Prevent the user from entering any leading whitespace
    const newVal = (e.target as HTMLInputElement).value.trimStart();
    const changed = inputVal.value !== newVal;

    if (changed) {
        inputVal.value = newVal;
        showSuggestions.value = true;
    } else {
        // query.value didn't change but the DOM element still has the whitespace.
        // Incrementing the key will force Vue to re-render with the cleaned value
        forceRenderKey.value++;
    }
}

// Emits an event with a selected tag or to trigger a post search
function onSubmit() {
    if (inputVal.value.length === 0) {
        emit("onSearch");
        return;
    }

    let tag: Tag;
    const negated = inputVal.value.startsWith("-");

    if (selectedIndex.value !== -1) {
        tag = suggestions.value[selectedIndex.value];
    } else {
        const value = negated ? inputVal.value.slice(1) : inputVal.value;
        const match = suggestions.value.find((t) => t.name === value);

        // User is submitting a raw tag that wasn't in the search suggestions
        if (match) {
            tag = match;
        } else {
            tag = {
                count: 0,
                name: value,
                type: "unknown",
            };
        }
    }

    inputVal.value = "";
    selectedIndex.value = -1;
    emit("onTagSelect", tag, negated);
}

// Setup a debounce to fetch search results shortly after the user stops typing
watch(inputVal, (query, _, onCleanup) => {
    onCleanup(() => clearTimeout(timer.value));

    selectedIndex.value = -1;

    if (query.length) {
        timer.value = setTimeout(() => doTagSearch(query), debounceMs);
    } else {
        suggestions.value = [];
    }
});
</script>

<template>
    <div
        class="container"
        @keydown.enter.prevent="onSubmit()"
        @keydown.tab.prevent="autoComplete()"
        @keydown.up.prevent="changeSelection(-1)"
        @keydown.down.prevent="changeSelection(1)"
        @keydown.esc.prevent="showSuggestions = false"
        ref="container"
    >
        <!-- forceRenderKey triggers a re-render when changed -->
        <template :key="forceRenderKey" />

        <input
            class="input"
            type="text"
            ref="input"
            placeholder="e.g: blue sky"
            :value="inputVal"
            @input="onInput"
            @focus="
                selectedIndex = -1;
                showSuggestions = true;
            "
            autofocus
        />
        <SearchSuggestions
            class="suggestions"
            :tags="suggestions"
            :selected-index="selectedIndex"
            @click="onSuggestionClick"
            v-if="showSuggestions"
        />
    </div>
    <button
        class="submit-btn btn-primary btn-rounded"
        type="submit"
        @click="onSubmit"
        :disabled="showSpinner"
    >
        <span v-if="!showSpinner">search</span>
        <span v-else class="spinner"></span>
    </button>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.input {
    background-color: #252525;
    border: 1px solid #555;
    color: #ddd;
    display: block;
    width: 100%;
    padding-left: 0.5rem;
    height: 40px;
}

.container {
    position: relative;
    margin-bottom: 0.8rem;
}

.suggestions {
    position: absolute;
    width: 100%;
    z-index: 1;
    padding-bottom: 1rem;
}

.submit-btn {
    $spinner-size: 20px;

    display: block;
    width: 100%;
    line-height: $spinner-size;

    span {
        display: block;
        margin: 0 auto;
    }

    .spinner {
        width: $spinner-size;
        height: $spinner-size;
        border-radius: 50%;
        border: 2px solid #fff;
        border-bottom-color: transparent;
        animation: linear 1s spin-anim infinite;
        display: block;

        @keyframes spin-anim {
            from {
                rotate: 0;
            }
            to {
                rotate: 360deg;
            }
        }
    }
}
</style>
