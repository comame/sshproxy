<script lang="ts">
  import type { ComponentType, SvelteComponent } from "svelte";
  import { onDestroy, tick } from "svelte";
  import { currentRoute } from "./routerImpl";

  let CurrentRouteComponent: ComponentType | null = null;

  let temporaryElementForRef: HTMLDivElement | null;
  let renderTarget: HTMLElement | null;
  let previousRendered: SvelteComponent | null = null;

  const unsubscribe = currentRoute.subscribe(async (v) => {
    if (v === null) {
      CurrentRouteComponent = null;
      return;
    }
    CurrentRouteComponent = v[0];

    // テンプレート部分で普通に <CurrentRouteComponent /> と書いても、再代入時に更新してくれないので、
    // やむを得ず手動でレンダリングさせる
    await tick();
    if (temporaryElementForRef) {
      renderTarget = temporaryElementForRef.parentElement;
      temporaryElementForRef.remove();
      temporaryElementForRef = null;
    }

    previousRendered?.$destroy();
    previousRendered = new CurrentRouteComponent({
      target: renderTarget!,
    });
  });

  onDestroy(() => {
    console.log("destroy router outlet");
    unsubscribe();
    previousRendered?.$destroy;
  });
</script>

{#if CurrentRouteComponent === null}
  <h1>Not Found</h1>
{/if}

<!-- 親要素の参照を取得するためのダミー要素 -->
<div bind:this={temporaryElementForRef}>Loading router...</div>
