<script setup lang="ts">
import ChangeLog from "./ChangeLog.vue";
import { changeLog } from "@/changelog";
import { COMMIT_SHA, LAST_COMMIT_DATE } from "@/config";
import { ref } from "vue";

const showChangelog = ref(false);
</script>

<template>
    <p>
        Booruview makes it easy to search and browse
        <a href="https://gelbooru.com" target="_blank">Gelbooru</a> posts. Check
        the help tab for more info.
    </p>
    <p>
        This project is open source and development is ongoing. Feedback and
        suggestions are welcome. You can use the
        <a href="https://github.com/jessienyan/booruview/issues" target="_blank"
            >Github issue tracker</a
        >
        or send me an
        <a href="mailto:216619670+jessienyan@users.noreply.github.com">email</a
        >.
    </p>

    <h4>
        <button class="btn-changelog" @click="showChangelog = !showChangelog">
            Changelog
            <i
                class="bi"
                :class="{
                    'bi-caret-down-fill': !showChangelog,
                    'bi-caret-up-fill': showChangelog,
                }"
            ></i>
        </button>
    </h4>
    <div v-if="showChangelog">
        <ChangeLog
            v-for="change in changeLog"
            :date="change.date"
            :changes="change.changes"
        />

        <p class="version">
            latest commit
            <a
                :href="`https://github.com/jessienyan/booruview/commits/${COMMIT_SHA}/`"
                target="_blank"
                ><code>{{ COMMIT_SHA }}</code></a
            >
            @ {{ LAST_COMMIT_DATE }}
        </p>
    </div>
</template>

<style scoped>
h4 {
    margin: 1rem 0;
    color: #ccc;
}

.btn-changelog {
    background: none;
    border: none;
    color: inherit;
    font-size: inherit;
    font-weight: inherit;
    cursor: pointer;
    padding: 0;
}

.version {
    font-family: "Courier New", Courier, monospace;
    font-size: 12px;
    text-align: right;
}
</style>
