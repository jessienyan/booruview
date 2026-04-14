<script setup lang="ts">
import { type RouteLocationRaw, RouterLink } from "vue-router";

const {
    currentPage,
    maxPage,
    totalCount,
    prevTo,
    nextTo,
    prevDisabled,
    nextDisabled,
    spinner,
} = defineProps<{
    currentPage: number;
    maxPage: number;
    totalCount: number;
    prevTo: RouteLocationRaw;
    nextTo: RouteLocationRaw;
    prevDisabled?: boolean;
    nextDisabled?: boolean;
    spinner?: "prev" | "next" | "none";
}>();

const fmt = new Intl.NumberFormat();
</script>

<template>
    <footer class="page-nav">
        <div class="nav-btns">
            <RouterLink v-if="currentPage > 1" :to="prevTo">
                <button
                    class="btn-primary btn-rounded"
                    :disabled="prevDisabled"
                    v-if="currentPage > 1"
                >
                    <div v-if="spinner === 'prev'" class="spinner">
                        <span class="spinner-inner"></span>
                    </div>
                    <template v-else>
                        <i class="bi bi-arrow-left"></i> prev
                    </template>
                </button>
            </RouterLink>
            <RouterLink v-if="currentPage < maxPage" :to="nextTo">
                <button
                    class="btn-primary btn-rounded"
                    :disabled="nextDisabled"
                >
                    <div v-if="spinner === 'next'" class="spinner">
                        <span class="spinner-inner"></span>
                    </div>
                    <template v-else>
                        next <i class="bi bi-arrow-right"></i>
                    </template>
                </button>
            </RouterLink>
        </div>
        <p>
            page {{ currentPage }} of {{ maxPage }} ({{
                fmt.format(totalCount)
            }}
            results)
        </p>
    </footer>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.page-nav {
    margin-top: 40px;
    text-align: center;
}

.nav-btns {
    display: flex;
    justify-content: center;
    gap: 20px;

    button {
        position: relative;
        width: 80px;
        height: 40px;
    }
}

.end-notice {
    color: orange;
}

.spinner {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);

    width: 20px;
    height: 20px;

    .spinner-inner {
        width: 100%;
        height: 100%;

        border-radius: 50%;
        border: 2px solid #fff;
        border-bottom-color: transparent;
        animation: linear 1s spin-anim infinite;
        display: block;

        @keyframes spin-anim {
            from {
                transform: rotate(0);
            }
            to {
                transform: rotate(360deg);
            }
        }
    }
}
</style>
