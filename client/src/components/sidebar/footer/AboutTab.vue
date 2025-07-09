<script setup lang="ts">
import ChangeLog from "./ChangeLog.vue";
import { changeLog } from "@/changelog";
import Collapsable from "@/components/Collapsable.vue";
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

    <Collapsable text="changelog">
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
    </Collapsable>
</template>

<style scoped>
.version {
    font-family: "Courier New", Courier, monospace;
    font-size: 12px;
    text-align: right;
}
</style>
