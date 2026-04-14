import type { Plugin } from 'vue'
import { nextTick } from 'vue'
import type {
  RouteLocationNormalized,
  RouteLocationNormalizedLoaded,
  Router,
} from 'vue-router'
import { isNavigationFailure } from 'vue-router'
import type {
  NavigationType,
  RouterScrollBehaviorOptions,
  ScrollPositionCoordinates,
  ScrollPositionCoordinatesGroup,
} from './types'

const STATE_KEY = 'vueRouterScroller'
const DEFAULT_INTERVAL = 200

/**
 * Setup router scroll behavior directly with a router instance.
 */
export function setupRouterScroller(
  router: Router,
  options: RouterScrollBehaviorOptions,
) {
  if (router.options.scrollBehavior) {
    console.warn(
      '`scrollBehavior` options in Vue Router is overwritten by `vue-router-scroller` plugin, you can remove it from createRouter()',
    )
  }

  router.options.scrollBehavior = () => {}
  options.storeInterval = options.storeInterval ?? DEFAULT_INTERVAL

  const positionsMap = new Map<string, ScrollPositionCoordinatesGroup>()
  let saveTimerId: number | null = null

  function startTimer(scrollKey: string) {
    if (saveTimerId !== null)
      return

    // Periodically save scroll positions
    saveTimerId = window.setInterval(() => {
      const pos = capturePositions(options)
      positionsMap.set(scrollKey, pos)
      history.replaceState({ ...history.state, [STATE_KEY]: pos }, '')
    }, options.storeInterval)
  }

  function stopTimer() {
    if (saveTimerId === null)
      return

    clearInterval(saveTimerId)
    saveTimerId = null
  }

  // Stop saving scroll positions while the state is being manipulated by
  // the browser/vue-router. Note that we can't listen to `popstate` because
  // vue-router also listens to it, and its listener is setup first. Listening
  // to `popstate` here would actually handle the event after all the callbacks
  // which is too late
  router.beforeEach(() => stopTimer())

  router.afterEach((to, from, failure) => {
    if (isNavigationFailure(failure))
      return

    const key = getScrollKey(to.fullPath)
    const pos = history.state[STATE_KEY] || positionsMap.get(key)
    const type = history.state[STATE_KEY] ? 'history' : 'push'

    nextTick(() => {
      applyPositions(to, from, pos, type, options)

      // Safe to start storing again
      startTimer(key)
    })
  })
}

/**
 * Set up router scroll behavior as a Vue plugin.
 *
 * @example
 * ```ts
 * import { createRouter } from 'vue-router'
 * import { createRouterScroller } from 'vue-router-better-scroller'
 *
 * const app = createApp(App)
 * const router = createRouter({ ... })
 *
 * app.use(router)
 * app.use(createRouterScroller({ ... }) // <-- this
 *
 * app.mount('#app')
 * ```
 */
export function createRouterScroller(
  options: RouterScrollBehaviorOptions,
): Plugin {
  return {
    install(app) {
      const router = app.config.globalProperties.$router
      if (!router) {
        throw new Error(
          'Router instance is not found on this Vue app. This plugin should be installed after Vue Router.',
        )
      }
      setupRouterScroller(router, options)
    },
  }
}

function getScrollKey(path: string, delta = 0): string {
  const position: number = history.state ? history.state.position - delta : -1
  return position + path
}

function capturePositions(options: RouterScrollBehaviorOptions) {
  const pos: ScrollPositionCoordinatesGroup = {}
  for (const [selector] of Object.entries(options.selectors)) {
    const element = querySelector(selector)
    if (!element)
      continue
    pos[selector] = getScrollPosition(element)
  }
  return pos
}

function querySelector(name: string) {
  if (typeof window === 'undefined')
    return undefined
  if (name === 'body')
    return document.body
  if (name === 'window')
    return window
  return document.querySelector(name)
}

function getScrollPosition(el: Element | Window): ScrollPositionCoordinates {
  if (el instanceof Window)
    return { left: window.scrollX, top: window.scrollY }
  else return { left: el.scrollLeft, top: el.scrollTop }
}

async function applyPositions(
  to: RouteLocationNormalized,
  from: RouteLocationNormalizedLoaded,
  pos: ScrollPositionCoordinatesGroup | undefined,
  type: NavigationType,
  options: RouterScrollBehaviorOptions,
) {
  for (const [selector, handler] of Object.entries(options.selectors)) {
    const element = querySelector(selector)
    if (!element)
      continue

    let position = pos?.[selector]
    if (typeof handler === 'function') {
      const result = await handler({
        to,
        from,
        element,
        selector,
        type,
        savedPosition: position,
      })
      if (!result)
        continue

      if (result !== true)
        position = result
    }
    else if (handler === true) {
      // by default, when navigate through a link, we don't preserve scroll position
      if (type === 'push')
        position = undefined
    }

    element.scrollTo({
      behavior: options.behavior,
      ...(position || { top: 0, left: 0 }),
    })
  }
}
