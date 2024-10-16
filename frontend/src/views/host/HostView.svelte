<script lang="ts">
  import HostQuizListView from "./HostQuizListView.svelte";
  import HostLobbyView from "./HostLobbyView.svelte";
  import HostPlayView from "./HostPlayView.svelte";
  import type { Quiz } from "../../model/quiz";
  import { HostGame, state } from "../../service/host/host";
  import { GameState } from "../../service/net";
  import HostIntermissionView from "./HostIntermissionView.svelte";

  let game = new HostGame();
  let active = false;

  function onHost(event: {detail: Quiz}) {
    game.hostQuiz(event.detail.id)
    active = true
  }

  let views: Record<GameState, any> = {
    [GameState.Lobby]: HostLobbyView,
    [GameState.Play]: HostPlayView,
    [GameState.Intermission]: HostIntermissionView,
    [GameState.Reveal]: HostPlayView,
  }
</script>

{#if active}
  <svelte:component this={views[$state]} {game}/>
{:else}
  <HostQuizListView on:host={onHost} />
{/if}