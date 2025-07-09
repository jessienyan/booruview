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
        Booruview is an open source <a href="https://gelbooru.com" target="_blank">Gelbooru</a> viewer. It's simple, fast, and works on desktop and mobile.
    </p>
    <ul class="links">
        <li>
            <a href="https://github.com/jessienyan/booruview" target="_blank">GitHub<i class="bi bi-github"></i></a>
        </li>
        <li>
            <a href="https://github.com/jessienyan/booruview/issues" target="_blank">Report a bug<i class="bi bi-bug-fill"></i></a>
        </li>
        <li>
            <a href="https://freesurveys.org/s/afOcVj9z-f" target="_blank">Give anonymous feedback<i class="bi bi-send-fill"></i></a>
        </li>
    </ul>

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
    font-size: 0.75rem;
    text-align: right;
}

.links .bi {
    margin-left: 0.4rem;
}
</style>
