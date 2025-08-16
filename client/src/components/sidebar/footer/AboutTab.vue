<script setup lang="ts">
import ChangeLog from "./ChangeLog.vue";
import { changeLog } from "@/changelog";
import Collapsable from "@/components/Collapsable.vue";
import { COMMIT_SHA, LAST_COMMIT_DATE } from "@/config";
import CodebergLogo from "@/assets/codeberg.svg";
</script>

<template>
    <p>
        Booruview is an open source
        <a href="https://gelbooru.com" target="_blank">Gelbooru</a> viewer. It's
        simple, fast, and works on desktop and mobile.
    </p>
    <ul class="links">
        <li>
            <a href="https://codeberg.org/jessienyan/booruview" target="_blank"
                ><img class="codeberg" :src="CodebergLogo" />Source code</a
            >
        </li>
        <li>
            <a
                href="https://github.com/jessienyan/booruview/issues"
                target="_blank"
                ><i class="bi bi-bug-fill"></i>Report a bug</a
            >
        </li>
        <li>
            <a href="https://freesurveys.org/s/afOcVj9z-f" target="_blank"
                ><i class="bi bi-send-fill"></i>Give anonymous feedback</a
            >
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

<style lang="scss" scoped>
@import "@/assets/colors";

.version {
    font-family: "Courier New", Courier, monospace;
    font-size: 0.75rem;
    text-align: right;
}

.links .bi,
.codeberg {
    margin-right: 0.4rem;
    vertical-align: middle;
}
</style>
