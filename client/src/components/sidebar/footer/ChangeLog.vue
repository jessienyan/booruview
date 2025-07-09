<script setup lang="ts">
import type { ChangeLogEntry } from "@/changelog";

const { date, changes } = defineProps<ChangeLogEntry>();
</script>

<template>
    <div class="entry">
        <h5>{{ date }}</h5>
        <ul>
            <li v-for="change in changes" v-html="change"></li>
        </ul>
    </div>
</template>
