<script setup lang="ts">
import { ref, Teleport } from "vue";
import ImportGelbooruFavs from "./ImportGelbooruFavs.vue";

const showFavsModal = ref(false);
</script>

<template>
    <div class="header-container">
        <button class="btn-primary btn-rounded" @click="showFavsModal = true">
            import gelbooru favs
        </button>

        <Teleport to="body">
            <ImportGelbooruFavs
                v-if="showFavsModal"
                @close="showFavsModal = false"
            />
        </Teleport>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.header-container {
    padding: 0.8em;
}
</style>
