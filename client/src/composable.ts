import {
	type ComputedRef,
	computed,
	inject,
	type MaybeRefOrGetter,
	onMounted,
	onUnmounted,
	readonly,
	ref,
	type ShallowRef,
	toValue,
} from "vue";

import store from "./store";
import type { RefOrGetter } from "./types";

export function useDismiss(el: MaybeRefOrGetter<MaybeRefOrGetter<HTMLElement | null>[]>, onDismiss: () => void) {
	function handler(e: MouseEvent) {
		const clickedOutside =
			toValue(el).findIndex(v => {
				const child = toValue(v);

				if (!child) {
					return false;
				}

				return child.contains(e.target as Node);
			}) === -1; // Clicked element was not found in any of `el`

		if (clickedOutside) {
			onDismiss();
		}
	}

	onMounted(() => document.addEventListener("click", handler));
	onUnmounted(() => document.removeEventListener("click", handler));
}

export function useIsVideo(post: MaybeRefOrGetter<Post>): ComputedRef<boolean> {
	return computed(() => {
		const val = toValue(post);

		return val.image_url.endsWith(".mp4") || val.image_url.endsWith(".webm");
	});
}

type StationaryClickReturn = {
	mouseDown: (e: MouseEvent) => void;
	mouseUp: (e: MouseEvent) => void;
};

// Fires a click event only if the cursor didn't move.
export function useStationaryClick(onClick: (e: MouseEvent) => void): StationaryClickReturn {
	const originX = ref(0);
	const originY = ref(0);

	// Number of pixels the cursor can move and still be considered a stationary click
	const allowedDistance = 10;

	function mouseDown(e: MouseEvent) {
		originX.value = e.x;
		originY.value = e.y;
	}

	function mouseUp(e: MouseEvent) {
		const dist = Math.sqrt((e.x - originX.value) ** 2 + (e.y - originY.value) ** 2);

		// Not a stationary click
		if (dist > allowedDistance) {
			return;
		}

		onClick(e);
	}

	return { mouseDown, mouseUp };
}

export function useDontShowAgain(id: string) {
	const flag = localStorage.getItem(id);
	const show = ref(flag === null);

	function ack() {
		localStorage.setItem(id, "1");
	}

	function onHide() {
		ack();
		show.value = false;
	}

	return { show, onHide, ack };
}

const now = ref(new Date());
setInterval(() => {
	now.value = new Date();
}, 5 * 1000);

/** Returns a reactive `new Date()` that periodically updates */
export function useDateNow() {
	return readonly(now);
}

/** Returns a reactive function that converts a Date into a relative string like "3 hours ago" */
export function useRelativeTime() {
	function timeString(date: Date) {
		return computed(() => {
			const now = useDateNow().value;
			const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);
			const minutes = Math.floor(seconds / 60);
			const hours = Math.floor(minutes / 60);
			const days = Math.floor(hours / 24);
			const months = Math.floor(days / 30);
			const years = Math.floor(days / 365);

			if (seconds < 60) {
				return "just now";
			} else if (minutes < 60) {
				return minutes === 1 ? "1 minute ago" : `${minutes} minutes ago`;
			} else if (hours < 24) {
				return hours === 1 ? "1 hour ago" : `${hours} hours ago`;
			} else if (days < 30) {
				return days === 1 ? "1 day ago" : `${days} days ago`;
			} else if (months < 12) {
				return months === 1 ? "1 month ago" : `${months} months ago`;
			} else {
				return years === 1 ? "1 year ago" : `${years} years ago`;
			}
		});
	}

	return timeString;
}

export function useMainContainer() {
	const mainContainer: Readonly<ShallowRef<HTMLElement>> = inject("mainContainer")!;
	return mainContainer;
}

export function useViewportSize() {
	const size = ref({ width: window.innerWidth, height: window.innerHeight });

	function updateSize() {
		size.value.width = window.innerWidth;
		size.value.height = window.innerHeight;
	}

	onMounted(() => window.addEventListener("resize", updateSize));
	onUnmounted(() => window.removeEventListener("resize", updateSize));

	return size;
}

// Rewrites an image URL to use the current CDN host
export function useGelbooruImageURL(url_: RefOrGetter<string>): ComputedRef<string> {
	return computed<string>(() => {
		const url = toValue(url_);

		if (store.cdnHosts === null) {
			return url;
		}

		if(store.cdnHosts.mediaProxy) {
			if(!url.includes(store.cdnHosts.image)) {
				// If it's missing the media proxy host, add it. Useful for favorites that are saved
				// with a static URL
				return store.cdnHosts.image + url;
			} else {
				// The API rewrites URLs to use the proxy, so there's nothing needed here
				return url;
			}
		}

		const newURL = new URL(url);
		newURL.host = store.cdnHosts.image;

		return newURL.toString();
	});
}

// Rewrites a video URL to use the current CDN host
export function useGelbooruVideoURL(url_: RefOrGetter<string>): ComputedRef<string> {
	return computed<string>(() => {
		let url = toValue(url_);

		// Fix incorrect cdn host. This logic is also in the API but it also needs to be
		// in the frontend to handle favorites
		url = url.replace("video-cdn3", "video-cdn4");

		if (store.cdnHosts === null) {
			return url;
		}

		if(store.cdnHosts.mediaProxy) {
			if(!url.includes(store.cdnHosts.video)) {
				// If it's missing the media proxy host, add it. Useful for favorites that are saved
				// with a static URL
				return store.cdnHosts.video + url;
			} else {
				// The API rewrites URLs to use the proxy, so there's nothing needed here
				return url;
			}
		}

		const newURL = new URL(url);
		newURL.host = store.cdnHosts.video;

		return newURL.toString();
	});
}
