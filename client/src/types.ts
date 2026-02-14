import type { MaybeRefOrGetter } from "vue";
import type {SearchHistory} from "@/store";

export interface ChipActions {
	// buttons
	blacklist?: boolean;
	edit?: boolean;
	favorite?: boolean;
	includeExcludeRemove?: boolean;
	openInNewTab?: boolean;

	// if true, clicking the chip doesn't open a menu
	static?: boolean;
}

export type RefOrGetter<T> = Exclude<MaybeRefOrGetter<T>, T>;

