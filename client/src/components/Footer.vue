<script setup lang="ts">
import { computed } from "vue";
import { RouterLink } from "vue-router";
import store from "@/store";

const fmt = new Intl.NumberFormat();
const currentPageText = computed(() => fmt.format(store.currentPage));
const maxPageText = computed(() => fmt.format(store.maxPage()));
const totalPostCountText = computed(() => fmt.format(store.totalPostCount));
</script>

<template>
    <footer class="page-nav">
        <p
            v-if="store.maxPage() > 200 && store.currentPage >= 200"
            class="end-notice"
        >
            Unfortunately, results past page 200 aren't viewable<br />because
            they are blocked by Gelbooru. :(
        </p>

        <div class="nav-btns">
            <RouterLink
                v-if="store.currentPage > 1"
                :to="{
                    name: 'search',
                    params: {
                        page: store.currentPage - 1,
                        query: $route.params.query,
                    },
                }"
            >
                <button
                    class="btn-primary btn-rounded"
                    :disabled="store.fetchingPosts"
                    v-if="store.currentPage > 1"
                >
                    <i class="bi bi-arrow-left"></i> prev
                </button>
            </RouterLink>
            <RouterLink
                v-if="store.currentPage < store.maxPage()"
                :to="{
                    name: 'search',
                    params: {
                        page: store.currentPage + 1,
                        query: $route.params.query,
                    },
                }"
            >
                <button
                    class="btn-primary btn-rounded"
                    :disabled="store.fetchingPosts || store.currentPage >= 200"
                >
                    next <i class="bi bi-arrow-right"></i>
                </button>
            </RouterLink>
        </div>
        <p>
            page {{ currentPageText }} of {{ maxPageText }} ({{
                totalPostCountText
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
