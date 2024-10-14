<script lang="ts">
  import { GameState } from "../../service/net";
  import { PlayerGame, state } from "../../service/player/player";
  import PlayerJoinView from "./PlayerJoinView.svelte";
  import PlayerLobbyView from "./PlayerLobbyView.svelte";
  import PlayerPlayView from "./PlayerPlayView.svelte";

  let game = new PlayerGame();
  let active = false

  function onJoin(){
    active = true
  }

  let views: Record<GameState, any> = {
    [GameState.Lobby]: PlayerLobbyView,
    [GameState.Play]: PlayerPlayView,
    [GameState.Reveal]: PlayerLobbyView
  }
</script>

{#if active}
  <svelte:component this={views[$state]} {game}} />
{:else}
  <PlayerJoinView on:join={onJoin} {game} />
{/if}