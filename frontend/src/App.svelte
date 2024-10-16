<script lang="ts">
  import PlayerView from "./views/player/PlayerView.svelte";
  import HostView from "./views/host/HostView.svelte";
  import Router from "svelte-spa-router";


  let routes = {
    '/': PlayerView,
    '/host': HostView
  }
</script>

<Router  {routes} />