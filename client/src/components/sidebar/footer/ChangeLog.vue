<script setup lang="ts">
import type { ChangeLogEntry } from "@/changelog";

const { date, changes } = defineProps<ChangeLogEntry>();
</script>

<template>
    <div class="entry">
        <p>
            <b>{{ date }}</b>
        </p>
        <ul>
            <li v-for="change in changes" v-html="change"></li>
        </ul>
    </div>
</template>

<style scoped>
p {
    color: #ccc;
}

p,
ul {
    margin: 10px 0;
}

.entry {
    margin-bottom: 20px;
}
</style>
