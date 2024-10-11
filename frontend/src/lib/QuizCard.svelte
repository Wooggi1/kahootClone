<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import Button from "./Button.svelte";
  import type { Quiz } from "../model/quiz";

  const dispatch = createEventDispatcher()

  export let quiz: Quiz
  function host(){
    dispatch("host", quiz)
  }
</script>

<div class="flex justify-between items-center bg-white border p-4 rounded-xl">
  <p>{quiz.name}</p>
  <Button on:click={host}>Host</Button>
</div>