<script setup lang="ts">
import type { ChangeLogEntry } from "@/changelog";

const { date, changes } = defineProps<ChangeLogEntry>();
</script>

<template>
    <h3>{{ date }}</h3>
    <ul>
        <li v-for="change in changes" v-html="change"></li>
    </ul>
</template>
