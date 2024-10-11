<script lang="ts">
  import { players, type HostGame } from "../../service/host/host";
  import Button from "../../lib/Button.svelte";
  import PlayerNameCard from "../../lib/lobby/PlayerNameCard.svelte";
  
  export let game: HostGame

  function start(){
    game.start()
  }
</script>

<div class="p-8 bg-purple-500 min-h-screen w-full">
  <div class="flex justify-end">
    <Button on:click={start}>Start Game</Button>
  </div>
  <div class="flex flex-wrap gap-2">
    {#each $players as player(player.id)}
      <PlayerNameCard {player} />
    {/each}
  </div>
</div>