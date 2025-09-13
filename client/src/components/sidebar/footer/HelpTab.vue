<script setup lang="ts">
import Chip from "@/components/tag-chip/Chip.vue";
import store from "@/store";
import { computed } from "vue";

const emit = defineEmits(["on-close"]);
const blueSkyTag = computed<TagChip>(() => {
    let style: ChipStyle = "default";
    if (store.query.isIncluded("blue_sky")) {
        style = "checkmark";
    } else if (store.query.isExcluded("blue_sky")) {
        style = "strikethrough";
    }
    return {
        style,
        tag: {
            count: 0,
            name: "blue_sky",
            type: "tag",
        },
    };
});
</script>

<template>
    <h2>shortcuts</h2>
    <ul>
        <li><i class="bi bi-phone"></i> Swipe left or right to change pages</li>
        <li><i class="bi bi-phone"></i> Pinch and zoom when viewing a post</li>
        <li>
            <i class="bi bi-display"></i> <code>Tab</code> when looking up a tag
            to autofill
        </li>
        <li>
            <i class="bi bi-display"></i> <code>Up</code> or
            <code>Down</code> to select a tag suggestion
        </li>
        <li>
            <i class="bi bi-display"></i> <code>Left</code> or
            <code>Right</code> to change posts
        </li>
        <li><i class="bi bi-display"></i> <code>F</code> to favorite a post</li>
    </ul>

    <h2>tips</h2>
    <ul>
        <li>
            Click the <i class="bi bi-info-circle"></i> icon when viewing a post
            to see its tags
        </li>
        <li>
            Any tag you find can be clicked to edit your search. No need to open
            the sidebar every time
        </li>
        <li>Blacklist tags you never want to see</li>
        <li>Favorite tags to find them quicker</li>
        <li>See your search history in the recent tab</li>
        <li>
            Tags are also links that will open a new page with just that tag as
            the search (ctrl+click, right click, or long press on mobile)
        </li>
        <li>
            Booruview is updated regularly, check the changelog for more info
        </li>
        <li>
            You can help make the site better by sending feedback in the about
            tab :)
        </li>
    </ul>

    <h2>searching</h2>
    <p>
        Like Gelbooru, Booruview uses tags as the building blocks for searching.
        Tags are buttons you can interact with by clicking or tapping them. Try
        it:
    </p>
    <p>
        <Chip :can-edit="false" :show-heart="false" :tag="blueSkyTag" />
    </p>
    <p>
        You can look up tags by typing them into the search box. Be sure to
        enter one tag at a time.
    </p>
    <p>
        Also, when looking up tags, don't worry about typing underscores. You'll
        get the same results whether you type:
    </p>
    <ul>
        <li><code>blue sky</code> or</li>
        <li><code>blue_sky</code></li>
    </ul>
    <p class="important">
        <i class="bi bi-info-circle"></i> The search box is <b>only</b> for
        looking up tags to add to your search. It's not for the search query
        itself.
    </p>

    <h3>example #1</h3>
    <p>
        On Gelbooru, if you wanted to search posts with one girl and no
        sunglasses, it would look like this:
    </p>
    <p>
        <code>1girl -sunglasses</code>
    </p>
    <p>
        On Booruview, you would first look up <code>1girl</code> and confirm it,
        then look up <code>-sunglasses</code> and confirm it.
    </p>
    <p>
        Alternatively, you could look up <code>sunglasses</code> (no hyphen)
        then exclude the tag by clicking it.
    </p>
    <p>Your search query should look like this:</p>
    <p>
        <Chip
            :can-interact="false"
            :can-edit="false"
            :tag="{
                style: 'default',
                tag: { count: 0, name: '1girl', type: 'tag' },
            }"
        />
        <Chip
            :can-interact="false"
            :can-edit="false"
            :tag="{
                style: 'strikethrough',
                tag: { count: 0, name: 'sunglasses', type: 'tag' },
            }"
        />
    </p>
    <p>Not this (blue tags are "raw"):</p>
    <p>
        <Chip
            :can-interact="false"
            :can-edit="false"
            :tag="{
                style: 'default',
                tag: { count: 0, name: '1girl -sunglasses', type: 'unknown' },
            }"
        />
    </p>

    <h3>example #2</h3>
    <p>
        One exception to the "lookup one tag at a time" rule is when doing an OR
        query. Booruview uses the same style as Gelbooru:
    </p>
    <p><code>{blue_sky ~ sunset}</code></p>
    <p>In Booruview it is the same:</p>
    <p>
        <Chip
            :can-interact="false"
            :can-edit="false"
            :tag="{
                style: 'default',
                tag: { count: 0, name: '{blue_sky ~ sunset}', type: 'unknown' },
            }"
        />
    </p>
    <p class="important">
        <i class="bi bi-info-circle"></i> Tags in an OR query (or any raw query)
        <b>must</b> use underscores instead of spaces.
        <code>{blue sky ~ sunset}</code> will not work.
    </p>

    <h3>search reference</h3>
    <ul>
        <li>Enter one tag at a time</li>
        <li>Searching without tags shows the most recent posts</li>
        <li><code>blue_sky</code> → search for blue_sky</li>
        <li><code>-blue_sky</code> → exclude blue_sky</li>
        <li><code>blue*</code> → tags that start with blue</li>
        <li><code>{blue_sky ~ sunset}</code> → blue_sky OR sunset</li>
        <li><code>sort:score</code> → sort by most popular</li>
        <li><code>rating:general</code> → SFW only</li>
    </ul>
    <p>
        For a complete list of all the different search filters, check the
        <a
            href="https://gelbooru.com/index.php?page=wiki&s=&s=view&id=26263"
            target="_blank"
            >Gelbooru search cheatsheet</a
        >.
    </p>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

code {
    font-family: "Courier New", Courier, monospace;
    background-color: rgba(121, 121, 121, 0.2);
    padding: 0 3px;
    border: 1px solid #4c4c4c;
    border-radius: 2px;
}

ul {
    padding-left: 25px;
    margin: 16px 0;
}

.important {
    background-color: $color-primary-darker;
    color: $color-primary-light;
    padding: 10px;
}

p :deep(.chip) {
    margin-right: 0.3rem;
}

h1,
h2,
h3 {
    margin-top: 25px;
}

p,
li {
    line-height: 120%;
}
</style>
