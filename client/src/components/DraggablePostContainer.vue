<script setup lang="ts">
import { computed, ref } from "vue";
import PostContainer from "./PostContainer.vue";

const {scrollContainer, posts} = defineProps<{
    scrollContainer: HTMLElement;
    posts: Post[];
}>();

const emit = defineEmits<{
    change: [posts: Post[]]
}>();

const draggingPostIndex = ref<number|null>(null);
const dropTargetIndex = ref<number|null>(null);

// "Working" view when dragging posts around
const postsWithDragOrder = computed(() => {
    if(draggingPostIndex.value == null || dropTargetIndex.value == null) {
        return posts;
    }

    const dragIndex = draggingPostIndex.value;
    const dropIndex = dropTargetIndex.value;

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
    ordered[dropIndex] = posts[draggingPostIndex.value];

    return ordered;
});

function onDragStart(e: DragEvent, index: number) {
    e.dataTransfer!.effectAllowed = "move";
    draggingPostIndex.value = index;
}

function onDragEnter(e: DragEvent, index: number) {
    if(draggingPostIndex.value == null) {
        return;
    }

    dropTargetIndex.value = index;
}

function onDragLeave(e: DragEvent, index: number) {
    if(draggingPostIndex.value == null) {
        return;
    }

    // dropTargetIndex.value = null;
}

function onDragEnd(e: DragEvent) {
    if(draggingPostIndex.value == null) {
        return;
    }

    if(dropTargetIndex.value != null) {
        e.preventDefault();
        // Emit the current working state as a way to commit the change
        emit("change", [...postsWithDragOrder.value]);
    }

    draggingPostIndex.value = null;
    dropTargetIndex.value = null;
}

const postDragId = computed(() => {
    if(draggingPostIndex.value == null) {
        return undefined;
    }

    return posts[draggingPostIndex.value].id;
})
</script>

<template>
    <PostContainer
        :keyed="true"
        :posts="postsWithDragOrder"
        :scroll-container="scrollContainer"
        :postDragId="postDragId"

        @post-dragstart="onDragStart"
        @post-dragenter="onDragEnter"
        @post-dragleave="onDragLeave"
        @post-dragend="onDragEnd"
    />
</template>
