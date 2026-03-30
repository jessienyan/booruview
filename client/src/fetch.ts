function delay(ms: number): Promise<void> {
	return new Promise(resolve => setTimeout(resolve, ms));
}

const defaultBackoffMS = 500;
const defaultRetryFn = (resp: Response) => resp.status >= 500;

type RetryOpts = {
	maxRetries: number;

	// The backoff time, in milliseconds. Total backoff time is linear and scales
	// with the number of attempts. The first retry happens immediately.
	// By default this is 500ms
	backoffMS?: number;

	// Retry function that should return true if the request should be retried.
	// By default, retries happen on 500+ errors.
	// Network errors are always retried.
	retry?: (resp: Response) => boolean;
}

// Fetches with an automatic backoff if the request fails
export async function fetchWithRetry(opts: RetryOpts, ...args: Parameters<typeof window.fetch>): Promise<Response> {
	opts.backoffMS = opts.backoffMS ?? defaultBackoffMS;
	opts.retry = opts.retry ?? defaultRetryFn;

	let remaining = opts.maxRetries;

	do {
		try {
			const resp = await fetch(...args);
			if(!opts.retry(resp)) {
				return resp;
			}
		} catch(e) {
			if(remaining === 0) {
				throw e;
			}
		}

		const attempt = opts.maxRetries - remaining;
		await delay(attempt * opts.backoffMS);

		remaining--;
	} while(remaining > 0);

	throw "shouldn't happen";
}
