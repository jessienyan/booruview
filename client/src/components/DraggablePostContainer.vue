<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from "vue";
import PostContent from "./PostContent.vue";
import store from "@/store";

type CroppedPost = {
    post: Post;
    cropped: boolean;
    renderHeight: number;
    index: number;
};

type ColumnDimensions = {
    count: number;
    width: number;
};

const { scrollContainer, posts } = defineProps<{
    scrollContainer: HTMLElement;
    posts: Post[];
}>();
const container = useTemplateRef("container");
const containerWidth = ref(0);

const maxPostHeight = computed(() => store.settings.maxPostHeight ?? 99999);
const colGap = 5;
const postGap = 5;

const theme = {
    colGap: colGap + "px",
    postGap: postGap + "px",
};

const draggingPost = ref<Post|null>(null);
const dropTarget = ref<Post|null>(null);

const orderedPosts = computed<CroppedPost[]>(() => {
    if(dropTarget.value == null || draggingPost.value == null || columnDimensions.value == null) {
        return [];
    }

    const dropIndex = posts.findIndex(p => p.id === dropTarget.value!.id);
    const dragIndex = posts.findIndex(p => p.id === draggingPost.value!.id);
    let ordered = [...posts];

    if(dragIndex < dropIndex) {
        // Shift right
        for(let i = dragIndex; i < dropIndex; i++) {
            ordered[i] = ordered[i+1]
        }
    } else {
        // Shift left
        for(let i = dragIndex; i > dropIndex; i--) {
            ordered[i] = ordered[i-1]
        }
    }

    // Update drop target
    ordered[dropIndex] = draggingPost.value;

    const columnWidth = columnDimensions.value.width;

    return ordered.map<CroppedPost>((p, i) => {
        const zoom = columnWidth / p.width;
        const renderHeight = p.height * zoom;
        const cropped = renderHeight > maxPostHeight.value;
        return {
            post: p,
            cropped,
            renderHeight: cropped ? maxPostHeight.value : renderHeight,
            index: i,
        };
    });
});

const columnDimensions = computed<ColumnDimensions | null>(() => {
    if (!container.value) {
        return null;
    }

    const ret: ColumnDimensions = { count: 1, width: 0 };

    if (store.settings.columnSizing === "fixed") {
        ret.count = store.settings.columnCount;
    } else {
        const colWithGap = store.settings.columnWidth + colGap;

        // prettier-ignore
        ret.count =
            1 +                                         // Always at least 1 column
            Math.max(0,
                Math.floor(
                    (containerWidth.value - store.settings.columnWidth) /    // First column is just the column width
                    colWithGap)                         // Remaining columns also have a gap
                );
    }

    ret.width = (containerWidth.value - (ret.count - 1) * colGap) / ret.count;

    return ret;
});

function onResize() {
    if (!container.value) {
        return;
    }

    containerWidth.value = container.value.clientWidth;
}

const postsByColumn = computed<CroppedPost[][]>(() => {
    if (columnDimensions.value == null) {
        return [];
    }

    const { count: columnCount } = columnDimensions.value;

    if (columnCount <= 1) {
        return [orderedPosts.value];
    }

    let ordered: CroppedPost[][] = [];
    let colHeight: number[] = [];
    for (let i = 0; i < columnCount; i++) {
        ordered = ordered.concat([[]]);
        colHeight = colHeight.concat(0);
    }

    for (const p of orderedPosts.value) {
        let shortestCol = 0;

        for (let i = 1; i < columnCount; i++) {
            if (colHeight[i] < colHeight[shortestCol]) {
                shortestCol = i;
            }
        }

        ordered[shortestCol] = ordered[shortestCol].concat(p);
        colHeight[shortestCol] += p.renderHeight + postGap;
    }

    return ordered;
});

function onDragStart(e: DragEvent, index: number) {
    e.dataTransfer!.effectAllowed = "move";
    draggingPost.value = orderedPosts.value[index].post;
}

function onDragEnter(e: DragEvent, index: number) {
    if(draggingPost.value === null || draggingPost.value.id === orderedPosts.value[index].post.id) {
        return;
    }

    dropTarget.value = orderedPosts.value[index].post;
}

function onDragLeave(e: DragEvent, index: number) {
    if(draggingPost.value === null || draggingPost.value.id === orderedPosts.value[index].post.id) {
        return;
    }

    dropTarget.value = null;
}

onMounted(() => {
    new ResizeObserver(onResize).observe(container.value!);
});
</script>

<template>
    <div class="post-container" ref="container">
        <div class="post-column" v-for="col in postsByColumn" >
            <PostContent
                v-for="post in col"
                :key="post.post.id"
                :post="post.post"
                :renderHeight="post.renderHeight"
                :maxHeight="maxPostHeight"
                :cropped="post.cropped"
                :scrollContainer="scrollContainer"

                @dragstart="(e: DragEvent) => onDragStart(e, post.index)"
                @dragenter="(e: DragEvent) => onDragEnter(e, post.index)"
                @dragleave="(e: DragEvent) => onDragLeave(e, post.index)"
            />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.post-container {
    display: flex;
    flex-direction: row;
    gap: v-bind("theme.colGap");
}

.post-column {
    display: flex;
    flex-direction: column;
    flex: 1;
    gap: v-bind("theme.postGap");

    &:last-of-type {
        margin-right: 0;
    }
}
</style>
