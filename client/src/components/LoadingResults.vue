<script setup lang="ts"></script>

<template>
    <div class="container"><span class="spinner"></span></div>
</template>

<style scoped lang="scss">
@import "@/assets/colors";

.container {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 100;
    display: flex;
    justify-content: center;
    align-items: center;
    color: white;
}

.spinner {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    border: 4px solid $color-primary;
    border-bottom-color: $color-primary-lighter;
    animation: linear 0.5s spin-anim infinite;
    display: block;

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
