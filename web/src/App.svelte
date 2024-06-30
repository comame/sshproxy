<script lang="ts">
  import Outlet from "./lib/router/Outlet.svelte";
  import Router from "./lib/router/Router.svelte";
  import Index from "./pages/Index.svelte";
  import NotFound from "./pages/NotFound.svelte";

  let value = "some text";
</script>

<input bind:value />
<div>value: {value}</div>

<Router
  routes={{
    "/": Index,
  }}
  fallback={NotFound}
  hashRouter
/>

<ul>
  <li><a href="#/">/</a></li>
  <li><a href="#/abc">/abc</a></li>
</ul>
<Outlet />
