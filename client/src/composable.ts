import {
    computed,
    onMounted,
    onUnmounted,
    readonly,
    ref,
    toValue,
    type ComputedRef,
    type MaybeRefOrGetter,
} from "vue";

export function useDismiss(
    el: MaybeRefOrGetter<HTMLElement | null>[],
    onDismiss: () => void,
) {
    function handler(e: MouseEvent) {
        const clickedOutside =
            el.findIndex((v) => {
                const real = toValue(v);
                if (real === null) {
                    return false;
                }

                return real.contains(e.target as Node);
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

        return (
            val.image_url.endsWith(".mp4") || val.image_url.endsWith(".webm")
        );
    });
}

type StationaryClickReturn = {
    mouseDown: (e: MouseEvent) => void;
    mouseUp: (e: MouseEvent) => void;
};

// Fires a click event only if the cursor didn't move.
export function useStationaryClick(
    onClick: (e: MouseEvent) => void,
): StationaryClickReturn {
    const originX = ref(0);
    const originY = ref(0);

    // Number of pixels the cursor can move and still be considered a stationary click
    const allowedDistance = 10;

    function mouseDown(e: MouseEvent) {
        originX.value = e.x;
        originY.value = e.y;
    }

    function mouseUp(e: MouseEvent) {
        const dist = Math.sqrt(
            Math.pow(e.x - originX.value, 2) + Math.pow(e.y - originY.value, 2),
        );

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

    function onHide() {
        localStorage.setItem(id, "1");
        show.value = false;
    }

    return { show, onHide };
}

export function useNewFeatureIndicator(id: string, until: Date) {
    // Feature indicator is no longer needed and will never be shown for new visitors
    if (new Date() > until) {
        const show = ref(false);
        return { show, onSeen: () => {} };
    }

    const { show, onHide: onSeen } = useDontShowAgain("feat-" + id);
    return { show, onSeen };
}

const now = ref(new Date());
setInterval(() => (now.value = new Date()), 5 * 1000);

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
                return minutes === 1
                    ? "1 minute ago"
                    : `${minutes} minutes ago`;
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
