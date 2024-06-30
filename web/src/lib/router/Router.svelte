<script lang="ts">
  import { type ComponentType } from "svelte";
  import { setupRouter, type routesDef } from "./routerImpl";

  export let routes: routesDef;
  export let hashRouter = false;
  export let fallback: ComponentType | null = null;
</script>

<svelte:window use:setupRouter={{ routes, hashRouter, fallback }} />
