<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from "vue";
import store from "@/store";
import PostContent from "./PostContent.vue";

type CroppedPost = {
	post: Post;
	cropped: boolean;
	renderHeight: number;
	index: number;
	key?: number;
};

type ColumnDimensions = {
	count: number;
	width: number;
};

const { scrollContainer, posts, keyed, postDragId } = defineProps<{
	scrollContainer: HTMLElement;
	posts: Post[];
	keyed: boolean;
	postDragId?: number;
}>();

const emit = defineEmits<{
	"post-dragstart": [event: DragEvent, postIndex: number];
	"post-dragenter": [event: DragEvent, postIndex: number];
	"post-dragleave": [event: DragEvent, postIndex: number];
	"post-dragend": [event: DragEvent, postIndex: number];
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
			1 + // Always at least 1 column
			Math.max(
				0,
				Math.floor(
					(containerWidth.value - store.settings.columnWidth) / // First column is just the column width
						colWithGap,
				), // Remaining columns also have a gap
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

const croppedPosts = computed<CroppedPost[]>(() => {
	if (columnDimensions.value == null) {
		return [];
	}

	const { width: columnWidth } = columnDimensions.value;

	return posts.map<CroppedPost>((post, index) => {
		const zoom = columnWidth / post.width;
		const renderHeight = post.height * zoom;
		const cropped = renderHeight > maxPostHeight.value;
		return {
			index,
			post,
			cropped,
			renderHeight: cropped ? maxPostHeight.value : renderHeight,
			key: keyed ? post.id : undefined,
		};
	});
});

const orderedPosts = computed<CroppedPost[][]>(() => {
	if (columnDimensions.value == null) {
		return [];
	}

	const { count: columnCount } = columnDimensions.value;

	if (columnCount <= 1) {
		return [croppedPosts.value];
	}

	let ordered: CroppedPost[][] = [];
	let colHeight: number[] = [];
	for (let i = 0; i < columnCount; i++) {
		ordered = ordered.concat([[]]);
		colHeight = colHeight.concat(0);
	}

	for (const p of croppedPosts.value) {
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

onMounted(() => {
	new ResizeObserver(onResize).observe(container.value!);
});
</script>

<template>
    <div class="post-container" ref="container">
        <div class="post-column" v-for="col in orderedPosts">
            <PostContent
                v-for="post in col"
                :post="post.post"
                :renderHeight="post.renderHeight"
                :maxHeight="maxPostHeight"
                :cropped="post.cropped"
                :scrollContainer="scrollContainer"
                :key="post.key"
                :beingDragged="postDragId === post.post.id"
                @dragstart="
                    (e: DragEvent) => emit('post-dragstart', e, post.index)
                "
                @dragenter="
                    (e: DragEvent) => emit('post-dragenter', e, post.index)
                "
                @dragleave="
                    (e: DragEvent) => emit('post-dragleave', e, post.index)
                "
                @dragend="(e: DragEvent) => emit('post-dragend', e, post.index)"
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
