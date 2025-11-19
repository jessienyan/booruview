<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useMainContainer } from "@/composable";

const top = ref("");
const scrollContainer = useMainContainer();
onMounted(() => (top.value = `${scrollContainer.value.scrollTop}px`));
</script>

<template>
    <div class="container" :style="{ top }">
        <span class="spinner"></span>
    </div>
</template>

<style scoped lang="scss">
@import "@/assets/colors";

.container {
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 3;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(2px);
}

.spinner {
    width: 70px;
    height: 70px;
    border-radius: 50%;
    border: 5px solid #705680;
    border-bottom-color: rgb(112, 86, 128);
    border-bottom-color: #8a709a;
    animation: linear 0.5s spin-anim infinite;
    display: block;
    filter: drop-shadow(0 0 5px black);

    @keyframes spin-anim {
        from {
            rotate: 0;
        }
        to {
            rotate: 360deg;
        }
    }
}
</style>
